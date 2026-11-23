package task

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

// gateProcessor lets a test control exactly when Process finishes.
// started gets a signal when Process begins; Process then blocks until release is closed.
type gateProcessor struct {
	started chan struct{}
	release chan struct{}
}

func newGate() gateProcessor {
	return gateProcessor{started: make(chan struct{}, 16), release: make(chan struct{})}
}

func (g gateProcessor) Process(input string) (*Result, error) {
	g.started <- struct{}{}
	<-g.release
	return &Result{Words: 1}, nil
}

// panicProcessor panics on the input "boom" and analyzes anything else normally.
type panicProcessor struct{}

func (panicProcessor) Process(input string) (*Result, error) {
	if input == "boom" {
		panic("kaboom")
	}
	return TextAnalyzer{}.Process(input)
}

func waitStarted(t *testing.T, ch <-chan struct{}) {
	t.Helper()
	select {
	case <-ch:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for a worker to start")
	}
}

func mustCreate(t *testing.T, s *Service, input string) int {
	t.Helper()
	task, err := s.CreateTask(input)
	if err != nil {
		t.Fatalf("CreateTask: %v", err)
	}
	return task.ID
}

func mustState(t *testing.T, s *Service, id int, want State) {
	t.Helper()
	got, err := s.GetTaskByID(id)
	if err != nil {
		t.Fatalf("GetTaskByID(%d): %v", id, err)
	}
	if got.State != want {
		t.Errorf("task %d state: got %v, want %v", id, got.State, want)
	}
}

// Two workers must be inside Process at the same moment: proves real concurrency.
func TestWorkersOverlap(t *testing.T) {
	g := newGate()
	s := NewService(g)
	s.Start(2, 4)

	id1 := mustCreate(t, s, "one")
	id2 := mustCreate(t, s, "two")
	if err := s.Submit(id1); err != nil {
		t.Fatal(err)
	}
	if err := s.Submit(id2); err != nil {
		t.Fatal(err)
	}

	waitStarted(t, g.started) // both signals arrive while neither task has finished
	waitStarted(t, g.started)

	mustState(t, s, id1, Processing)
	mustState(t, s, id2, Processing)

	close(g.release)
	s.Shutdown()

	mustState(t, s, id1, Completed)
	mustState(t, s, id2, Completed)
}

// One worker is busy and the queue (size 1) holds one task: the next Submit must be rejected.
func TestQueueFull(t *testing.T) {
	g := newGate()
	s := NewService(g)
	s.Start(1, 1)

	id1 := mustCreate(t, s, "a")
	id2 := mustCreate(t, s, "b")
	id3 := mustCreate(t, s, "c")

	if err := s.Submit(id1); err != nil {
		t.Fatal(err)
	}
	waitStarted(t, g.started) // worker took task 1, queue is empty again

	if err := s.Submit(id2); err != nil { // fills the one queue slot
		t.Fatal(err)
	}

	err := s.Submit(id3) // no room left
	if !errors.Is(err, ErrQueueFull) {
		t.Fatalf("got %v, want ErrQueueFull", err)
	}
	mustState(t, s, id3, Pending) // rolled back, not stuck in Queued

	close(g.release)
	s.Shutdown()
	mustState(t, s, id1, Completed)
	mustState(t, s, id2, Completed)
	mustState(t, s, id3, Pending)
}

func TestSubmitBeforeStartAndAfterShutdown(t *testing.T) {
	s := NewService(TextAnalyzer{})
	id := mustCreate(t, s, "hello")

	if err := s.Submit(id); !errors.Is(err, ErrNotStarted) {
		t.Errorf("before Start: got %v, want ErrNotStarted", err)
	}

	s.Start(2, 4)
	s.Shutdown()

	if err := s.Submit(id); !errors.Is(err, ErrStopped) {
		t.Errorf("after Shutdown: got %v, want ErrStopped", err)
	}
	mustState(t, s, id, Pending)
}

func TestShutdownIsIdempotent(t *testing.T) {
	s := NewService(TextAnalyzer{})
	s.Shutdown() // never started: must not panic or hang
	s.Start(2, 4)
	s.Shutdown()
	s.Shutdown() // second call: must not panic on close
}

// Shutdown must wait for everything that was already queued.
func TestShutdownDrainsQueue(t *testing.T) {
	s := NewService(TextAnalyzer{})
	s.Start(3, 20)

	var ids []int
	for i := 0; i < 20; i++ {
		id := mustCreate(t, s, fmt.Sprintf("task number %d", i))
		if err := s.Submit(id); err != nil {
			t.Fatal(err)
		}
		ids = append(ids, id)
	}

	s.Shutdown()

	for _, id := range ids {
		mustState(t, s, id, Completed)
	}
}

// A panicking Processor must fail only its own task, and the worker must survive.
func TestProcessorPanicIsContained(t *testing.T) {
	s := NewService(panicProcessor{})
	s.Start(1, 4) // ONE worker: if it died, the second task would never finish

	bad := mustCreate(t, s, "boom")
	good := mustCreate(t, s, "fine words here")
	if err := s.Submit(bad); err != nil {
		t.Fatal(err)
	}
	if err := s.Submit(good); err != nil {
		t.Fatal(err)
	}

	s.Shutdown()

	mustState(t, s, bad, Failed)
	mustState(t, s, good, Completed)

	got, _ := s.GetTaskByID(bad)
	if got.Err == "" {
		t.Error("failed task should record an error message")
	}
}

// Many goroutines hammering the service at once. Run with -race: any missing lock fails here.
func TestConcurrentUse(t *testing.T) {
	s := NewService(TextAnalyzer{})
	s.Start(4, 100)

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			task, err := s.CreateTask(fmt.Sprintf("goroutine %d says hello", n))
			if err != nil {
				t.Errorf("CreateTask: %v", err)
				return
			}
			if err := s.Submit(task.ID); err != nil {
				t.Errorf("Submit: %v", err)
				return
			}
			s.GetTaskByID(task.ID) // concurrent readers too
		}(i)
	}
	wg.Wait()
	s.Shutdown()

	for id := 1; id <= 50; id++ {
		mustState(t, s, id, Completed)
	}
}

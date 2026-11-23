package task

import (
	"fmt"
	"strings"
	"sync"
)

type Service struct {
	mu      sync.RWMutex
	tasks   map[int]*Task
	nextID  int
	process Processor
	queue   chan int
	wg      sync.WaitGroup
	started bool
	stopped bool
}

func NewService(process Processor) *Service {
	return &Service{
		tasks:   make(map[int]*Task),
		nextID:  1,
		process: process,
	}
}

func (s *Service) CreateTask(input string) (*Task, error) {

	if strings.TrimSpace(input) == "" {
		return nil, ErrEmptyInput
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	task := &Task{ID: s.nextID, Input: input, State: Pending}
	s.tasks[task.ID] = task
	s.nextID++

	return task.clone(), nil
}

func (s *Service) GetTaskByID(id int) (*Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, err := s.get(id)
	if err != nil {
		return nil, err
	}

	return task.clone(), nil

}

func (s *Service) get(id int) (*Task, error) {
	task, ok := s.tasks[id]
	if !ok {
		return nil, fmt.Errorf("task %d: %w", id, ErrNotFound)
	}

	return task, nil
}

func (s *Service) ProcessTask(id int) (*Task, error) {

	s.mu.Lock()
	task, err := s.get(id)
	if err != nil {
		s.mu.Unlock()
		return nil, err
	}

	if err := task.transition(Processing); err != nil {
		s.mu.Unlock()
		return nil, err
	}

	input := task.Input
	s.mu.Unlock()

	result, procErr := s.safeProcess(input)

	s.mu.Lock()
	defer s.mu.Unlock()

	if procErr != nil {
		task.State = Failed
		task.Err = procErr.Error()
		return nil, fmt.Errorf("task %d: processing failed: %w", id, procErr)
	}

	task.Result = result
	if err := task.transition(Completed); err != nil {
		return nil, err
	}

	return task.clone(), nil

}

func (s *Service) Start(workers, queueSize int) {

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.started {
		return
	}

	s.queue = make(chan int, queueSize)
	s.started = true

	for i := 0; i < workers; i++ {
		s.wg.Add(1)
		go s.worker(s.queue)
	}
}

func (s *Service) worker(queue <-chan int) {
	defer s.wg.Done()

	for id := range queue {
		s.ProcessTask(id)
	}

}

func (s *Service) Submit(id int) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.started {
		return ErrNotStarted
	}

	if s.stopped {
		return ErrStopped
	}

	task, err := s.get(id)
	if err != nil {
		return err
	}

	if err := task.transition(Queued); err != nil {
		return err
	}

	select {
	case s.queue <- id:
		return nil
	default:
		task.State = Pending
		return fmt.Errorf("task %d: %w", id, ErrQueueFull)
	}

}

// safeProcess runs the Processor and turns a panic into an ordinary error,
// so one bad task can never kill a worker goroutine.
func (s *Service) safeProcess(input string) (result *Result, err error) {
	defer func() {
		if r := recover(); r != nil {
			result = nil
			err = fmt.Errorf("processor panic: %v", r)
		}
	}()

	return s.process.Process(input)
}

// Shutdown stops accepting new tasks, lets the workers finish everything
// already queued, and returns when all workers have exited.
// It is safe to call more than once.
func (s *Service) Shutdown() {
	s.mu.Lock()
	if !s.started {
		s.mu.Unlock()
		return
	}
	if !s.stopped {
		s.stopped = true
		close(s.queue) // same lock as Submit, so nobody can send on a closed channel
	}
	s.mu.Unlock()

	// Wait must happen WITHOUT the lock: workers need it inside ProcessTask
	// to finish their current task, so holding it here would deadlock.
	s.wg.Wait()
}

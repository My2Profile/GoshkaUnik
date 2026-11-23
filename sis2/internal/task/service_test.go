package task

import (
	"errors"
	"testing"
)

func TestCreateEmpty(t *testing.T) {
	s := NewService(TextAnalyzer{})

	task, err := s.CreateTask("")

	if !errors.Is(err, ErrEmptyInput) {
		t.Errorf("got %v, want ErrEmptyInput", err)
	}
	if task != nil {
		t.Errorf("got %v, want nil", task)
	}
}

func TestCreateInvalidInput(t *testing.T) {
	cases := []struct {
		input string
	}{
		{""},
		{"	"},
		{"\n\t"},
		{"    "},
	}

	for _, tc := range cases {
		service := NewService(TextAnalyzer{})
		_, err := service.CreateTask(tc.input)
		if !errors.Is(err, ErrEmptyInput) {
			t.Errorf("got %v, want ErrEmptyInput", err)
		}
	}

}

func TestCreateAndGet(t *testing.T) {
	s := NewService(TextAnalyzer{})

	task, err := s.CreateTask("hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := s.GetTaskByID(task.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Get returns a copy, so compare the contents, not the pointers.
	if *got != *task {
		t.Errorf("got %v, want %v", got, task)
	}
}

func TestGetReturnsCopy(t *testing.T) {
	s := NewService(TextAnalyzer{})
	task, _ := s.CreateTask("hello")

	task.State = Failed // mutate our copy

	got, _ := s.GetTaskByID(task.ID)
	if got.State != Pending {
		t.Errorf("service state changed through a copy: got %v, want %v", got.State, Pending)
	}
}

func TestGetNotFound(t *testing.T) {
	s := NewService(TextAnalyzer{})

	_, err := s.GetTaskByID(5)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("got %v, want ErrNotFound", err)
	}
}

func TestProcessTask(t *testing.T) {
	s := NewService(TextAnalyzer{})

	task, err := s.CreateTask("Hello Go world!")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	done, err := s.ProcessTask(task.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if done.State != Completed {
		t.Errorf("state: got %v, want %v", done.State, Completed)
	}

	want := Result{Words: 3, Runes: 15, Lines: 1, UniqueWords: 3}
	if done.Result == nil || *done.Result != want {
		t.Errorf("result: got %v, want %v", done.Result, want)
	}
}

func TestProcessTwice(t *testing.T) {
	s := NewService(TextAnalyzer{})
	task, _ := s.CreateTask("Hello Go world!")

	if _, err := s.ProcessTask(task.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err := s.ProcessTask(task.ID)
	if !errors.Is(err, ErrInvalidTransition) {
		t.Errorf("got %v, want ErrInvalidTransition", err)
	}

	got, _ := s.GetTaskByID(task.ID) // fetch a fresh snapshot
	if got.State != Completed {
		t.Errorf("state: got %v, want %v", got.State, Completed)
	}
}

func TestIndependentTasks(t *testing.T) {
	s := NewService(TextAnalyzer{})
	t1, _ := s.CreateTask("one two")
	t2, _ := s.CreateTask("alpha beta gamma")

	if _, err := s.ProcessTask(t1.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got2, _ := s.GetTaskByID(t2.ID)
	if got2.State != Pending || got2.Result != nil {
		t.Errorf("task2 should be untouched, got state %v result %v", got2.State, got2.Result)
	}

	if _, err := s.ProcessTask(t2.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got1, _ := s.GetTaskByID(t1.ID)
	got2, _ = s.GetTaskByID(t2.ID)
	if got1.Result.Words != 2 || got2.Result.Words != 3 {
		t.Errorf("words: task1 %d (want 2), task2 %d (want 3)", got1.Result.Words, got2.Result.Words)
	}
}

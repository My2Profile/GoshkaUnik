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

	if got != task {
		t.Errorf("got %v, want %v", got, task)
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

	processedTask, err := s.ProcessTask(1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if processedTask.State != Completed {
		t.Errorf("got %v, want %v", processedTask.State, Completed)
	}

	if task.Result.Words != 3 {
		t.Errorf("got %v, want %v", task.Result.Words, 3)
	}

	if task.Result.Runes != 15 {
		t.Errorf("got %v, want %v", task.Result.Runes, 15)
	}

	if task.Result.Lines != 1 {
		t.Errorf("got %v, want %v", task.Result.Lines, 1)
	}

	if task.Result.UniqueWords != 3 {
		t.Errorf("got %v, want %v", task.Result.UniqueWords, 3)
	}

}

func TestProcessTwice(t *testing.T) {
	s := NewService(TextAnalyzer{})

	task, err := s.CreateTask("Hello Go world!")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = s.ProcessTask(task.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = s.ProcessTask(task.ID)
	if !errors.Is(err, ErrInvalidTransition) {
		t.Errorf("got %v, want ErrInvalidTransition", err)
	}

	if task.State != Completed {
		t.Errorf("got %v, want %v", task.State, Completed)
	}

	if task.Result.Words != 3 {
		t.Errorf("got %v, want %v", task.Result.Words, 3)
	}

	if task.Result.Runes != 15 {
		t.Errorf("got %v, want %v", task.Result.Runes, 15)
	}

	if task.Result.Lines != 1 {
		t.Errorf("got %v, want %v", task.Result.Lines, 1)
	}

	if task.Result.UniqueWords != 3 {
		t.Errorf("got %v, want %v", task.Result.UniqueWords, 3)
	}

}

func TestIndependentTasks(t *testing.T) {

	s := NewService(TextAnalyzer{})

	task1, err := s.CreateTask("one two")

	task2, err := s.CreateTask("alpha beta gamma")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = s.ProcessTask(task1.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = s.ProcessTask(task2.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if task1.State != Completed {
		t.Errorf("got %v, want %v", task1.State, Completed)
	}

	if task2.State != Completed && task2.Result != nil {
		t.Errorf("got %v, want %v", task2.State, Completed)
	}

}

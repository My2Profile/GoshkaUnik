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

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	task2, err := s.CreateTask("alpha beta gamma")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = s.ProcessTask(task1.ID)
	if task2.State != Pending { // which state is a task that was never processed?
		t.Errorf("task2 state: got %v, want %v", task2.State, Pending)
	}
	if task2.Result != nil { // no result yet: what value means "nothing"?
		t.Errorf("task2 result: got %v, want nil", task2.Result)
	}

	// 3. task1 has its own result
	if task1.Result.Words != 2 { // "one two" has how many words?
		t.Errorf("task1 words: got %v, want %v", task1.Result.Words, 2)
	}

	// 4. now process task2
	_, err = s.ProcessTask(task2.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if task2.Result.Words != 3 { // "alpha beta gamma"
		t.Errorf("task2 words: got %v, want %v", task2.Result.Words, 3)
	}

	// 5. and task1 must NOT have changed
	if task1.Result.Words != 2 {
		t.Errorf("task1 words changed: got %v, want %v", task1.Result.Words, 2)
	}

}

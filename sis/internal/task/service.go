package task

import (
	"fmt"
	"strings"
)

type Service struct {
	tasks   map[int]*Task
	nextID  int
	process Processor
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

	task := &Task{
		ID:     s.nextID,
		Input:  input,
		State:  Pending,
		Result: nil,
	}

	s.tasks[s.nextID] = task
	s.nextID++

	return task, nil
}

func (s *Service) GetTaskByID(id int) (*Task, error) {

	task, ok := s.tasks[id]
	if !ok {
		return nil, fmt.Errorf("task %d: %w", id, ErrNotFound)
	}

	return task, nil
}

func (s *Service) ProcessTask(id int) (*Task, error) {

	task, err := s.GetTaskByID(id)
	if err != nil {
		return nil, err
	}

	if err := task.transition(Processing); err != nil {
		return nil, err
	}

	task.Result, err = s.process.Process(task.Input)
	if err != nil {
		task.State = Pending
		return nil, fmt.Errorf("task %d: processing failed: %w", id, err)
	}

	if err := task.transition(Completed); err != nil {
		return nil, err
	}

	return task, nil

}

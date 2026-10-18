package task

import (
	"errors"
)

var (
	ErrEmptyInput = errors.New("input is empty or whitespace only")

	ErrNotFound = errors.New("task not found with that ID")

	ErrInvalidTransition = errors.New("invalid transition for task")
)

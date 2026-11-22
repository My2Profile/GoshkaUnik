package task

import (
	"errors"
)

var ErrNotFound = errors.New("task not found")

var ErrEmptyInput = errors.New("the input is empty")

var ErrInvalidTransition = errors.New("transition is invalid")

var ErrQueueFull = errors.New("queue is Full")

var ErrStopped = errors.New("submit after shutdown")

var ErrNotStarted = errors.New("submit before start")

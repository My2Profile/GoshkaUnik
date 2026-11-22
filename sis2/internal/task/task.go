package task

import (
	"fmt"
)

type State int

const (
	Pending State = iota
	Queued
	Processing
	Completed
	Failed
)

var stateNames = [...]string{"Pending", "Queued", "Processing", "Completed", "Failed"}

func (s State) String() string {
	if int(s) < 0 || int(s) >= len(stateNames) {
		return fmt.Sprintf("Unknown State (%d)", int(s))
	}

	return stateNames[s]
}

type Result struct {
	Words, Runes, Lines, UniqueWords int
}

type Task struct {
	ID     int
	Input  string
	State  State
	Result *Result
	Err    string
}

var allowed = map[State]map[State]bool{
	Pending:    {Queued: true, Processing: true},
	Queued:     {Processing: true, Pending: true},
	Processing: {Completed: true, Failed: true},
}

func (t *Task) transition(to State) error {

	if !allowed[t.State][to] {
		return fmt.Errorf("%w: %s -> %s", ErrInvalidTransition, t.State, to)
	}

	t.State = to
	return nil

}

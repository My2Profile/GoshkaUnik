package task

import (
	"fmt"
)

type State int

const (
	Pending State = iota
	Processing
	Completed
)

func (s State) String() string {
	switch s {
	case Pending:
		return "Pending"

	case Processing:
		return "Processing"
	case Completed:
		return "Completed"

	default:
		return fmt.Sprintf("Unknown State (%d)", s)
	}
}

type Result struct {
	Words       int
	Runes       int
	Lines       int
	UniqueWords int
}

type Task struct {
	ID     int
	Input  string
	State  State
	Result *Result
}

var allowed = map[State]State{
	Pending:    Processing,
	Processing: Completed,
}

func (t *Task) transition(to State) error {

	next, ok := allowed[t.State]
	if !ok || next != to {
		return fmt.Errorf("%w: %s -> %s", ErrInvalidTransition, t.State, to)
	}

	t.State = to
	return nil

}

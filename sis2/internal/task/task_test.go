package task

import (
	"errors"
	"testing"
)

func TestTransition(t *testing.T) {
	tests := []struct {
		name     string
		from, to State
		wantErr  bool
	}{
		{"pending to queued", Pending, Queued, false},
		{"completed to processing", Completed, Processing, true},
		{"queued to processing", Queued, Processing, false},
		{"pending to processing", Pending, Processing, false},
		{"processing to completed", Processing, Completed, false},
		{"processing to pending", Processing, Pending, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			task := &Task{State: tc.from}
			err := task.transition(tc.to)

			if tc.wantErr {
				if !errors.Is(err, ErrInvalidTransition) {
					t.Errorf("err = %v, want ErrInvalidTransition", err)
				}
				if task.State != tc.from {
					t.Errorf("state changed to %s on a failed transition", task.State)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if task.State != tc.to {
				t.Errorf("state = %s, want %s", task.State, tc.to)
			}

		})
	}
}

func TestStateString(t *testing.T) {

	tests := []struct {
		state State
		want  string
	}{
		{Pending, "Pending"},
		{Failed, "Failed"},
		{State(5), "Unknown State (5)"},
		{State(-1), "Unknown State (-1)"},
	}

	for _, tc := range tests {
		got := tc.state.String()
		if got != tc.want {
			t.Errorf("State(%d).String() = %q, want %q", int(tc.state), got, tc.want)
		}
	}

}

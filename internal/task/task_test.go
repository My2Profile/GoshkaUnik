package task

import "testing"

func TestTransition(t *testing.T) {
	cases := []struct {
		name    string
		from    State
		to      State
		wantErr bool
	}{
		{"pending to processing", Pending, Processing, true},
		{"processing to completed", Processing, Completed, true},
		{"pending to completed", Pending, Completed, false},
		{"completed to pending", Completed, Pending, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			task := Task{State: tc.from}
			err := task.transition(tc.to)
			if (err != nil) != tc.wantErr {
				t.Errorf("got err %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

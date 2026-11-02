package main

import (
	"errors"
	"fmt"
)

// Sentinel error for testing errors.Is propagation
var ErrSentinel = errors.New("sentinel configuration error")

type Task struct {
	Name string
	Run  func() error
}

// SafeRun executes all tasks, catching panics and wrapping errors.
func SafeRun(tasks []Task) []error {
	var errs []error

	for _, task := range tasks {
		err := executeTask(task)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

// executeTask runs a single task inside a panic recovery boundary.
func executeTask(task Task) (errOut error) {
	defer func() {
		if p := recover(); p != nil {
			// Convert panic to a formatted error
			errOut = fmt.Errorf("task %s: panic: %v", task.Name, p)
		}
	}()

	if err := task.Run(); err != nil {
		// Wrap standard errors with %w to preserve sentinel checks
		return fmt.Errorf("task %s: %w", task.Name, err)
	}

	return nil
}

func main() {
	tasks := []Task{
		{
			Name: "successful-task",
			Run: func() error {
				return nil // Succeeds
			},
		},
		{
			Name: "sentinel-task",
			Run: func() error {
				return ErrSentinel // Returns sentinel error
			},
		},
		{
			Name: "index-out-of-range",
			Run: func() error {
				s := []int{1, 2}
				_ = s[10] // Triggers runtime panic
				return nil
			},
		},
		{
			Name: "nil-map-write",
			Run: func() error {
				var m map[string]int
				m["key"] = 42 // Triggers runtime panic (assignment to entry in nil map)
				return nil
			},
		},
		{
			Name: "explicit-panic-task",
			Run: func() error {
				panic("config file is missing") // Explicit panic call
			},
		},
	}

	errs := SafeRun(tasks)

	fmt.Println("=== Individual Task Errors ===")
	for _, e := range errs {
		fmt.Println(e)
	}

	// Join all errors together
	combinedErr := errors.Join(errs...)

	fmt.Println("\n=== Verifying Sentinel with errors.Is ===")
	if errors.Is(combinedErr, ErrSentinel) {
		fmt.Println("Success: ErrSentinel was successfully found inside errors.Join(errs...)!")
	} else {
		fmt.Println("Failure: ErrSentinel could not be found.")
	}
}

package main

import (
	"fmt"
)

type ValidationError struct {
	Field string
}

func (e *ValidationError) Error() string {
	return "invalid field: " + e.Field
}

// Fixed validate function: returns a literal nil when there is no error
func validate(age int) error {
	if age < 0 {
		return &ValidationError{Field: "age"}
	}
	return nil // Clean literal nil return
}

func main() {
	err := validate(20) // age is valid
	fmt.Printf("err == nil: %v, Type: %T\n", err == nil, err)

	if err != nil {
		fmt.Println("Error caught:", err)
	} else {
		fmt.Println("Validation passed successfully!")
	}
}

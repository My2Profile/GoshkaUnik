package main

import (
	"errors"
	"fmt"
)

// 1. Define a sentinel error for missing students
var ErrNotFound = errors.New("student not found in gradebook")

// 2. Define a custom error type for invalid grades (must implement the error interface)
type GradeError struct {
	Name  string
	Grade int
}

// Implementing the error interface by adding Error() string method
func (e *GradeError) Error() string {
	return fmt.Sprintf("invalid grade %d for %s: must be between 0 and 100", e.Grade, e.Name)
}

// 3. Define the GradeBook struct wrapping a map
type GradeBook struct {
	records map[string]int
}

// Helper to initialize a GradeBook
func NewGradeBook() *GradeBook {
	return &GradeBook{
		records: make(map[string]int),
	}
}

// 4. Implement Get method
func (gb *GradeBook) Get(name string) (int, error) {
	grade, exists := gb.records[name]
	if !exists {
		return 0, ErrNotFound // Return sentinel error
	}
	return grade, nil
}

// 5. Implement Set method
func (gb *GradeBook) Set(name string, grade int) error {
	if grade < 0 || grade > 100 {
		return &GradeError{Name: name, Grade: grade} // Return custom pointer error
	}
	gb.records[name] = grade
	return nil
}

func addBonus(gb *GradeBook, name string, bonus int) error {
	if _, err := gb.Get(name); err != nil {
		return fmt.Errorf("cannot add bonus: %w", err) // Wrap the error
	}
	
	if grade, _ := gb.Get(name); grade+bonus > 100 {
		return &GradeError{Name: name, Grade: grade + bonus} // Return custom pointer error
	}

	gb.Set(name, bonus)
	return nil
	
}


func main() {
	gb := NewGradeBook()

	// Setting a valid grade
	err := gb.Set("Alice", 95)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Trying to set an invalid grade (triggers *GradeError)
	err = gb.Set("Bob", 105)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Getting an existing student
	if grade, err := gb.Get("Alice"); err == nil {
		fmt.Printf("Alice's grade: %d\n", grade)
	}

	// Getting a non-existent student (triggers ErrNotFound)
	if _, err := gb.Get("Charlie"); errors.Is(err, ErrNotFound) {
		fmt.Println("Charlie is not in the gradebook!")
	}
}

package main

import (
	"fmt"
)

type Resource struct {
	Name string
}

func Open(name string) (*Resource, error) {
	fmt.Println("  -> open:", name)
	return &Resource{Name: name}, nil
}

func (r *Resource) Close() {
	fmt.Println("  <- close:", r.Name)
}

// BAD: Defers pile up until the function exits completely
func processAllBad(names []string) {
	fmt.Println("--- processAllBad START ---")
	for _, name := range names {
		r, _ := Open(name)
		defer r.Close() // Queued up, closes happen AFTER the loop finishes!
	}
	fmt.Println("--- processAllBad Loop Finished (closes happen next) ---")
}

// GOOD: Helper function limits defer scope to each iteration
func processAllGood(names []string) {
	fmt.Println("--- processAllGood START ---")
	for _, name := range names {
		// Immediately invoked function expression (IIFE) scopes the defer
		func(n string) {
			r, _ := Open(n)
			defer r.Close() // Runs at the end of this iteration
		}(name)
	}
	fmt.Println("--- processAllGood END ---")
}

func main() {
	names := []string{"file1.txt", "file2.txt", "file3.txt"}

	processAllBad(names)
	fmt.Println()
	processAllGood(names)
}

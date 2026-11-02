package main

import (
	"fmt"
	"time"
)

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("[TRACE] %s took %v\n", name, time.Since(start))
	}
}

func heavyCalculation() {
	defer timer("heavyCalculation")() // Evaluates timer immediately, defers the returned closure

	time.Sleep(100 * time.Millisecond) // Simulate work
}

func main() {
	heavyCalculation()
}

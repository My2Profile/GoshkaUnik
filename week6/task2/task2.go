package main

import (
	"fmt"
	"time"
)

func memorize(f func(int) int) func(int) int {
	found := make(map[int]int)

	return func(n int) int {
		if result, ok := found[n]; ok {
			fmt.Printf("(memorized) %d -> %d\n", n, result)
			return result
		}

		fmt.Printf(" Not memorized: %d\n", n)
		fmt.Println("Calling the underlying function...")
		time.Sleep(100 * time.Millisecond)
		result := f(n)

		found[n] = result
		return result
	}
}

func main() {

	check := 0
	slowFunction := func(n int) int {
		check++
		return n * 2
	}

	memorizedFunction := memorize(slowFunction)
	fmt.Println("--- Call 1: Input 5 ---")
	fmt.Println("Result:", memorizedFunction(5))

	fmt.Println("\n--- Call 2: Input 5 (Repeated) ---")
	fmt.Println("Result:", memorizedFunction(5))

	fmt.Println("\n--- Call 3: Input 7 (New) ---")
	fmt.Println("Result:", memorizedFunction(7))

	fmt.Println("\n--- Call 4: Input 7 (Repeated) ---")
	fmt.Println("Result:", memorizedFunction(7))

	fmt.Printf("\n=== Proof ===\n")
	fmt.Printf("Total unique inputs: 2 (5 and 7)\n")
	fmt.Printf("Total actual underlying function calls: %d\n", check)

}

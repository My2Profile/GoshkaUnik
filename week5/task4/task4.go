package main

import (
	"cmp"
	"fmt"
)

// Define custom Number constraint using tilde for types like Score
type Number interface {
	~int | ~int64 | ~float64
}

// 1. Sum and Average
func Sum[T Number](nums []T) T {
	var s T
	for _, n := range nums {
		s += n
	}
	return s
}

func Average[T Number](nums []T) (float64, bool) {
	if len(nums) == 0 {
		return 0, false
	}
	total := Sum(nums)
	// Convert generic sum to float64 for division
	return float64(total) / float64(len(nums)), true
}

// 2. MaxOf using cmp.Ordered
func MaxOf[T cmp.Ordered](items []T) (T, bool) {
	var zero T
	if len(items) == 0 {
		return zero, false
	}
	m := items[0]
	for _, v := range items[1:] {
		if v > m {
			m = v
		}
	}
	return m, true
}

// 3. IndexOf using comparable
func IndexOf[T comparable](items []T, target T) int {
	for i, v := range items {
		if v == target {
			return i
		}
	}
	return -1
}

// 4. Filter using any
func Filter[T any](items []T, keep func(T) bool) []T {
	var result []T
	for _, v := range items {
		if keep(v) {
			result = append(result, v)
		}
	}
	return result
}

type Score int

func main() {
	scores := []Score{85, 60, 92, 78}

	total := Sum(scores)
	avg, _ := Average(scores)
	fmt.Printf("Total Score: %d, Average: %.2f\n", total, avg)

	// Filter scores >= 75
	passing := Filter(scores, func(s Score) bool {
		return s >= 75
	})
	fmt.Println("Passing scores:", passing)
}

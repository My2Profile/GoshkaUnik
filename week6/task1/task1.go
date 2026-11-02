package main

import (
	"fmt"
)

func makeAccumulator() func(int) int {
	count := 0

	return func(n int) int {
		count += n
		return count
	}

}

func main() {

	acc := makeAccumulator()
	fmt.Println(acc(3))
	fmt.Println(acc(5))
	fmt.Println(acc(10)) // Output: 18

}

package main

import "fmt"

func main() {

	s1 := []int{1, 2, 3, 4, 5}

	s2 := make([]int, 3, 5)

	arr := []int{10, 20, 30, 40, 50}

	s3 := arr[1:4]

	var s4 []int

	s4 = append(s4, 1, 2)

	fmt.Println(s1, s2, s3, s4)

}

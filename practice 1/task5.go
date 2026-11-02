package main

import "fmt"

func modify(s []int) {
	s[0] = 99;
}

func main() {

	var arr []int

	for i := 0; i < 5; i++ {
		arr = append(arr, i)
	}

	fmt.Println(arr)
	
	modify(arr)

	fmt.Println(arr)



}

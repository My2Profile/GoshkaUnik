package main

import "fmt"

func main() {

	var arr [5][5]int

	for i := range arr {
		for j := range arr[i] {
			arr[i][j] = i*5 + j + 1
		}
	}

	fmt.Println(arr)

}

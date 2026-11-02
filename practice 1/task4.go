package main

import "fmt"

func main() {

	arr := [5]int{1, 2, 3, 4, 5}

	fmt.Println("Array before mods", arr)

	slice := arr[1:4]

	slice[0] = 99

	fmt.Println("Array after slice modification", arr)
	fmt.Println("Slice view", slice)

}

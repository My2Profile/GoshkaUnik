package main

import (
	"fmt"
)

func main() {
	x := 1
	defer fmt.Println("deferred x =", x)
	defer func() { fmt.Println("closure x =", x) }()
	x = 2
	fmt.Println("current x =", x)
}

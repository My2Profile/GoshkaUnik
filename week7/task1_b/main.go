package main

import (
	"fmt"
)

func catch(name string) {
	if r := recover(); r != nil {					// trying to recover?
		fmt.Println(name, "recovered:", r)
	}
}

func normalPanic() {
	defer catch("normal Panic")

	var m map[string]int
	m["x"] = 1										// Panic writing to the nil map
	fmt.Println("never printed?")
}

func forever(n int) int { return forever(n+1) + 1 }  // Forever loop

func overflow() {
	defer catch("overflow")
	forever(0)										// Running forever loop
}

func main() {
	normalPanic()
	fmt.Println("after normal panic the program should still be alive")

	overflow()
	fmt.Println("cannot see this line")				// the overflow is not a panic cannot recover it

}


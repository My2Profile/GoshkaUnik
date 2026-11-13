// Task 4.2: preemption. Run: go run ./task4/preempt   and   GODEBUG=asyncpreemptoff=1 go run ./task4/preempt
package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(1)
	go func() {
		for { // busy loop with no function calls
		}
	}()
	time.Sleep(10 * time.Millisecond)
	fmt.Println("main is still alive")
}

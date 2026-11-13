// Week 8, Task 2: First goroutines.
// Usage: go run . <part>
//
//	say      - part 1: say with sync.WaitGroup
//	broken   - part 2: the "two classic mistakes" program (run it 3 times and with -race)
//	atomic   - part 2: fix #1 with atomic.Int64
//	slice    - part 2: fix #2, each goroutine writes its own []int element
//	panic    - part 3: panic in goroutine killed the program (no recover in the goroutine)
//	recover  - part 3: defer recover INSIDE the goroutine keeps the program alive
//	workers  - participation task: "worker 5 is printed first"
package main

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

// ---------- Part 1 ----------

func say(msg string) {
	for i := 0; i < 3; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(msg, i)
	}
}

func partSay() {
	var wg sync.WaitGroup
	wg.Add(1) // Add BEFORE go
	go func() {
		defer wg.Done()
		say("goroutine")
	}()
	say("main")
	wg.Wait() // main waits, so both loops always finish
}

// ---------- Part 2 ----------

func partBroken() {
	var wg sync.WaitGroup
	counter := 0
	for i := 0; i < 1000; i++ {
		go func() {
			wg.Add(1) // WRONG: Wait() may run before this line
			defer wg.Done()
			counter++ // WRONG: data race
		}()
	}
	wg.Wait()
	fmt.Println("counter:", counter)
}

func partAtomic() {
	var wg sync.WaitGroup
	var counter atomic.Int64
	for i := 0; i < 1000; i++ {
		wg.Add(1) // before go
		go func() {
			defer wg.Done()
			counter.Add(1) // atomic read-modify-write
		}()
	}
	wg.Wait()
	fmt.Println("counter:", counter.Load())
}

func partSlice() {
	var wg sync.WaitGroup
	results := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results[i] = 1 // each goroutine owns its own element: no race
		}()
	}
	wg.Wait()
	sum := 0
	for _, v := range results {
		sum += v
	}
	fmt.Println("counter:", sum)
}

// ---------- Part 3 ----------

func partPanic() {
	// recover here does NOTHING for the other goroutine
	defer func() { fmt.Println("main recovered:", recover()) }()
	go func() {
		var m map[string]int
		m["x"] = 1 // panic in goroutine -> whole program crashes
	}()
	time.Sleep(100 * time.Millisecond)
	fmt.Println("never printed")
}

func partRecover() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		// The defer must be INSIDE the goroutine that may panic.
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("goroutine recovered:", r)
			}
		}()
		var m map[string]int
		m["x"] = 1
	}()
	wg.Wait()
	fmt.Println("main is still alive")
}

// ---------- Participation ----------

func partWorkers() {
	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println("worker", i, "done")
		}()
	}
	wg.Wait()
	fmt.Println("all workers finished")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run . say|broken|atomic|slice|panic|recover|workers")
		return
	}
	switch os.Args[1] {
	case "say":
		partSay()
	case "broken":
		partBroken()
	case "atomic":
		partAtomic()
	case "slice":
		partSlice()
	case "panic":
		partPanic()
	case "recover":
		partRecover()
	case "workers":
		partWorkers()
	default:
		fmt.Println("unknown part")
	}
}

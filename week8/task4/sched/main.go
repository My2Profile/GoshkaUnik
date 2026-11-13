// Task 4.1: 8 busy goroutines. Run: GODEBUG=schedtrace=500 go run ./task4/sched   and with GOMAXPROCS=1
package main

import (
	"sync"
	"time"
)

func busy(d time.Duration) {
	end := time.Now().Add(d)
	for time.Now().Before(end) {
	}
}

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			busy(2 * time.Second)
		}()
	}
	wg.Wait()
}

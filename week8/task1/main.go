// Task 1: Concurrency vs parallelism. Run: go run ./task1   and   GOMAXPROCS=1 go run ./task1
package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for d := 2; d < n; d++ { // deliberately slow
		if n%d == 0 {
			return false
		}
	}
	return true
}

func countPrimes(from, to int) int {
	c := 0
	for n := from; n < to; n++ {
		if isPrime(n) {
			c++
		}
	}
	return c
}

func cpuWork(workers int) (int, time.Duration) {
	const N = 300_000
	start := time.Now()
	results := make([]int, workers)
	var wg sync.WaitGroup
	chunk := N / workers
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results[w] = countPrimes(w*chunk, (w+1)*chunk)
		}()
	}
	wg.Wait()
	total := 0
	for _, r := range results {
		total += r
	}
	return total, time.Since(start)
}

func ioWork(workers int) time.Duration {
	start := time.Now()
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(100 * time.Millisecond)
		}()
	}
	wg.Wait()
	return time.Since(start)
}

func main() {
	cores := runtime.NumCPU()
	fmt.Printf("NumCPU=%d GOMAXPROCS=%d\n", cores, runtime.GOMAXPROCS(0))
	for _, w := range []int{1, cores, 2 * cores, 4 * cores} {
		p, d := cpuWork(w)
		fmt.Printf("CPU work, %3d goroutine(s): %d primes in %v\n", w, p, d.Round(time.Millisecond))
	}
	for _, w := range []int{1, 10, 1000} {
		fmt.Printf("I/O work, %4d goroutine(s): %v\n", w, ioWork(w).Round(time.Millisecond))
	}
}

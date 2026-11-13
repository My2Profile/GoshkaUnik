// Task 3: cost of a goroutine. Run: go run ./task3 [N]
package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

func main() {
	n := 100_000
	if len(os.Args) > 1 {
		n, _ = strconv.Atoi(os.Args[1])
	}
	var before, after runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&before)

	stop := make(chan struct{})
	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-stop
		}()
	}
	elapsed := time.Since(start)
	time.Sleep(200 * time.Millisecond)
	runtime.ReadMemStats(&after)

	fmt.Println("goroutines alive:", runtime.NumGoroutine())
	fmt.Printf("time to start %d goroutines: %v\n", n, elapsed)
	d := after.StackSys - before.StackSys
	fmt.Printf("stack memory: %.0f MB (%.2f KB per goroutine)\n", float64(d)/1e6, float64(d)/1024/float64(n))
	close(stop)
	wg.Wait()
	time.Sleep(100 * time.Millisecond)
	fmt.Println("after close:", runtime.NumGoroutine(), "goroutine(s)")
}

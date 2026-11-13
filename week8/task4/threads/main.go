// Task 4.3: threads created by time.Sleep vs blocking syscall.
package main

import (
	"fmt"
	"runtime/pprof"
	"sync"
	"syscall"
	"time"
)

func run(name string, n int, f func()) {
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); f() }()
	}
	wg.Wait()
	fmt.Printf("%-24s threads created so far: %d\n", name, pprof.Lookup("threadcreate").Count())
}

func main() {
	fmt.Printf("%-24s threads created so far: %d\n", "start", pprof.Lookup("threadcreate").Count())
	run("100 x time.Sleep", 100, func() { time.Sleep(200 * time.Millisecond) })
	run("100 x blocking syscall", 100, func() {
		ts := syscall.NsecToTimespec(int64(200 * time.Millisecond))
		syscall.Nanosleep(&ts, nil) // the whole OS thread is blocked
	})
}

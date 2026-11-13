// Task 5. Run: go run ./task5 [leak|profile|fixed]
package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

func report(step string) {
	fmt.Printf("%-28s goroutines: %d\n", step, runtime.NumGoroutine())
}

func main() {
	mode := "leak"
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}
	switch mode {
	case "leak": // 1. count the leaks
		report("start")
		for i := 0; i < 100; i++ {
			fetchLeaky(time.Millisecond) // timeout 1 ms < 50 ms query
		}
		time.Sleep(100 * time.Millisecond)
		report("after 100 fetchLeaky")
		for i := 0; i < 10; i++ {
			startTickerLeaky()
		}
		report("after 10 leaky tickers")
	case "profile": // 2. find them
		for i := 0; i < 100; i++ {
			fetchLeaky(time.Millisecond)
		}
		for i := 0; i < 10; i++ {
			startTickerLeaky()
		}
		time.Sleep(100 * time.Millisecond)
		fmt.Println("===== debug=1 =====")
		pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
		fmt.Println("===== debug=2 (first 25 lines) =====")
		var buf limited
		pprof.Lookup("goroutine").WriteTo(&buf, 2)
	case "fixed": // 3. prove the fix
		report("start")
		for i := 0; i < 100; i++ {
			fetchFixed(time.Millisecond)
		}
		time.Sleep(100 * time.Millisecond)
		report("after 100 fetchFixed")
		ctx, cancel := context.WithCancel(context.Background())
		for i := 0; i < 10; i++ {
			startTicker(ctx)
		}
		report("after 10 stoppable tickers")
		cancel()
		time.Sleep(100 * time.Millisecond)
		report("after cancel()")
	}
}

// limited prints only the first lines of the debug=2 dump.
type limited struct{ lines int }

func (l *limited) Write(p []byte) (int, error) {
	for _, b := range p {
		if l.lines >= 25 {
			break
		}
		fmt.Print(string(b))
		if b == '\n' {
			l.lines++
		}
	}
	return len(p), nil
}

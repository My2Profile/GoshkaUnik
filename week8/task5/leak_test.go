package main

import (
	"os"
	"runtime"
	"testing"
	"time"
)

// fetch is the function under test. By default it is the fixed version;
// run with LEAKY=1 to see the test FAIL for fetchLeaky.
func fetch(timeout time.Duration) (int, error) {
	if os.Getenv("LEAKY") == "1" {
		return fetchLeaky(timeout)
	}
	return fetchFixed(timeout)
}

func TestFetchDoesNotLeak(t *testing.T) {
	before := runtime.NumGoroutine()
	for i := 0; i < 100; i++ {
		fetch(time.Millisecond) // always times out (query takes 50 ms)
	}
	// give goroutines a short time to finish (slowQuery = 50 ms)
	time.Sleep(200 * time.Millisecond)
	after := runtime.NumGoroutine()
	if after > before {
		t.Fatalf("goroutine leak: before=%d after=%d (%d leaked)", before, after, after-before)
	}
}

package main

import (
	"sync"
	"testing"
)

var sink int

//go:noinline
func work() int { sink++; return sink }

func BenchmarkGoroutineStartStop(b *testing.B) {
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); work() }()
		wg.Wait()
	}
}

func BenchmarkFuncCall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		work()
	}
}

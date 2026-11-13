package main

import (
	"context"
	"errors"
	"time"
)

func slowQuery() int {
	time.Sleep(50 * time.Millisecond) // pretend to ask a database
	return 42
}

// LEAK 1: unbuffered channel. After a timeout nobody reads ch,
// so the goroutine blocks forever on "ch <- ..." (state: chan send).
func fetchLeaky(timeout time.Duration) (int, error) {
	ch := make(chan int)
	go func() {
		ch <- slowQuery()
	}()
	select {
	case v := <-ch:
		return v, nil
	case <-time.After(timeout):
		return 0, errors.New("timeout")
	}
}

// FIX 1: buffer of one -> the send always succeeds and the goroutine finishes.
func fetchFixed(timeout time.Duration) (int, error) {
	ch := make(chan int, 1)
	go func() {
		ch <- slowQuery()
	}()
	select {
	case v := <-ch:
		return v, nil
	case <-time.After(timeout):
		return 0, errors.New("timeout")
	}
}

// LEAK 2: a worker nobody can stop (state: sleep).
func startTickerLeaky() {
	go func() {
		for {
			time.Sleep(10 * time.Millisecond)
		}
	}()
}

// FIX 2: the caller can say "stop" through the context.
func startTicker(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(10 * time.Millisecond):
			}
		}
	}()
}

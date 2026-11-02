package main

import (
	"fmt"
	"time"
	"runtime"
)

type Record struct {
	ID   int
	Data [256]byte
}

var live []*Record
var sum int

func consume(r *Record) int { return r.ID + int(r.Data[0]) }

func main() {
	start := time.Now()
	for i := 0; i < 80_000; i++ {
		live = append(live, &Record{ID: i})
	}
	for i := 0; i < 10_000_000; i++ {
		r := &Record{ID: i} // garbage after one use
		sum += consume(r)
	}
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Println("GC cycles:", m.NumGC, "| heap from OS (MB):", m.HeapSys/1024/1024,
		"| total pause:", time.Duration(m.PauseTotalNs), "| time:", time.Since(start))
	runtime.KeepAlive(live)
}
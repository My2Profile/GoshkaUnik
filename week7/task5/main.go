package main

import (
	"fmt"
	"strings"
	"bytes"
	"runtime"
)

var kept [][]byte

func keepHeadersSlice(buf []byte) { kept = append(kept, buf[:16]) }
func keepHeadersClone(buf []byte) { kept = append(kept, bytes.Clone(buf[:16])) }


func heapMB() uint64 {
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.HeapAlloc / 1024 / 1024
}

func main() {
	for i := 0; i < 100; i++ {
		buf := make([]byte, 1<<20)
		keepHeadersSlice(buf)
	}
	fmt.Println(heapMB())

	kept = nil
	fmt.Println(heapMB())

	for j := 0; j < 100; j++ {
		buf := make([]byte, 1<<20)
		keepHeadersClone(buf)
	}
	
	fmt.Println("slices cloned:", heapMB(), "MB")

	var keptStr []string
	for i := 0; i < 100; i++ {
		big := strings.Repeat("x", 1<<20)
		keptStr = append(keptStr, big[:10])
	}
	fmt.Println("substrings sliced:", heapMB(), "MB")

	keptStr = nil
	for i := 0; i < 100; i++ {
		big := strings.Repeat("x", 1<<20)
		keptStr = append(keptStr, strings.Clone(big[:10]))
	}
	fmt.Println("substrings cloned:", heapMB(), "MB")

	runtime.KeepAlive(kept)
	runtime.KeepAlive(keptStr)
}

package main

import (
	"fmt"
	"strings"
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
		keepHeadersSlice()
	}
	fmt.Println(headMB())

	kept = nil
	fmt.Println(headMB())

	for j := 0; j < 100; j++ {
		buf := make([]byte, 1<<20)
		keepHeadersClone()
	}
	fmt.Println(headMB())

	big := strings.Repeat("x", 1<<20)
	small := big[:10]

	runtime.KeepAlive(kept)

}

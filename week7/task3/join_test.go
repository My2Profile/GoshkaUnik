package main

import "testing"

var ids = func() []int {
	// var ints [1000]int;
	var ints = make([]int, 1000)
	for i := 0; i < 1000; i++ {
		ints[i] = i
	}
	return ints
}()

func TestAllocs(t *testing.T) {
	slow := testing.AllocsPerRun(100, func() { JoinIDsSlow(ids) })
	fast := testing.AllocsPerRun(100, func() { JoinIDsFast(ids) })
	t.Logf("slow=%v fast=%v", slow, fast)
}

func BenchmarkSlow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		JoinIDsSlow(ids)
	}
}
func BenchmarkFast(b *testing.B) {
	for i := 0; i < b.N; i++ {
		JoinIDsFast(ids)
	}
}

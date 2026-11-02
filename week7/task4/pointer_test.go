package main

import "testing"

const n = 1000

var sinkP Point
var sinkPP *Point
var sinkS []int

func TestCollect(t *testing.T) {
	got := CollectWithCap(3)
	if len(got) != 3 || got[2] != 2 {
		t.Errorf("got %v, want [0 1 2]", got)
	}
}

func TestAllocs(t *testing.T) {
	nocap := testing.AllocsPerRun(100, func() { sinkS = CollectNoCap(n) })
	withcap := testing.AllocsPerRun(100, func() { sinkS = CollectWithCap(n) })
	t.Logf("nocap lol=%v withcap notlol=%v", nocap, withcap)
}

func BenchmarkNewPointValue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkP = NewPointValue(i, i)
	}
}

func BenchmarkNewPointPointer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkPP = NewPointPointer(i, i)
	}
}

func BenchmarkCollectNoCap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkS = CollectNoCap(n)
	}
}

func BenchmarkCollectWithCap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkS = CollectWithCap(n)
	}
}

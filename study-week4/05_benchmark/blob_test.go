package benchmark

import "testing"

// sink is package-level so the compiler cannot discard the call results
// (dead-code elimination) and quietly delete the work we want to measure.
var sink float64

// Anatomy of a benchmark:
//   - lives in a *_test.go file, named BenchmarkXxx, takes *testing.B
//   - setup goes BEFORE the loop
//   - the framework picks b.N (grows it until timing is stable)
func BenchmarkBlobValue(b *testing.B) {
	x := Blob{}
	for i := 0; i < b.N; i++ {
		sink = x.FirstValue()
	}
}

func BenchmarkBlobPointer(b *testing.B) {
	x := Blob{}
	for i := 0; i < b.N; i++ {
		sink = x.FirstPointer()
	}
}

func BenchmarkVecValue(b *testing.B) {
	v := Vec{1, 2}
	for i := 0; i < b.N; i++ {
		sink = v.SumValue()
	}
}

func BenchmarkVecPointer(b *testing.B) {
	v := Vec{1, 2}
	for i := 0; i < b.N; i++ {
		sink = v.SumPointer()
	}
}

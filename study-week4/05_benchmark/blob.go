// STUDY FILE (analog of Task 5: value vs pointer receiver benchmark).
// Types here are Blob (large) and Vec (small). Run the benchmarks with:
//
//	cd study-week4/05_benchmark
//	go test -bench=. -benchmem
//
// WHAT TO LEARN
//   - A value receiver COPIES the struct into the callee's stack frame on every
//     call. For a 16 KB struct that is a 16 KB memcpy per call.
//   - B/op and allocs/op count HEAP allocations only. The stack copy is not one,
//     so they are 0 for every variant. The cost shows up in the ns/op column.
//   - //go:noinline forbids the compiler from pasting the tiny method body into
//     the caller. If it inlined, it could skip the copy and the benchmark would
//     measure nothing.
package benchmark

type Blob struct {
	Label string
	Data  [2048]float64 // 2048 * 8 bytes = 16 KB
}

//go:noinline
func (b Blob) FirstValue() float64 { return b.Data[0] }

//go:noinline
func (b *Blob) FirstPointer() float64 { return b.Data[0] }

type Vec struct{ X, Y float64 } // 16 bytes: copying is trivial

//go:noinline
func (v Vec) SumValue() float64 { return v.X + v.Y }

//go:noinline
func (v *Vec) SumPointer() float64 { return v.X + v.Y }

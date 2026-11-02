package main

import "testing"

type Report struct {
	Title  string
	Values [1024]float64 // 8 KB payload (1024 * 8 bytes)
}

type Point struct {
	X, Y float64 // 16 bytes payload
}

//go:noinline
func (r Report) FirstValue() float64 { return r.Values[0] }

//go:noinline
func (r *Report) FirstPointer() float64 { return r.Values[0] }

//go:noinline
func (p Point) FirstValue() float64 { return p.X }

//go:noinline
func (p *Point) FirstPointer() float64 { return p.X }

// --- Benchmarks for Large Struct (Report - ~8 KB) ---

func BenchmarkReportValue(b *testing.B) {
	r := Report{}
	var result float64
	for i := 0; i < b.N; i++ {
		result = r.FirstValue()
	}
	_ = result
}

func BenchmarkReportPointer(b *testing.B) {
	r := Report{}
	var result float64
	for i := 0; i < b.N; i++ {
		result = r.FirstPointer()
	}
	_ = result
}

// --- Benchmarks for Small Struct (Point - 16 bytes) ---

func BenchmarkPointValue(b *testing.B) {
	p := Point{X: 10, Y: 20}
	var result float64
	for i := 0; i < b.N; i++ {
		result = p.FirstValue()
	}
	_ = result
}

func BenchmarkPointPointer(b *testing.B) {
	p := Point{X: 10, Y: 20}
	var result float64
	for i := 0; i < b.N; i++ {
		result = p.FirstPointer()
	}
	_ = result
}

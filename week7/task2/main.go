package main

import "fmt"

type Point struct{ X, Y int }

var global *Point 																//heap

func f1() int 				{ x := 42; return x }   							// stack
func f2() *int 				{ x := 42; return &x }  							// heap returns pointer
func f3(p *Point) int 		{ return p.X + p.Y }								// stack?
func f4()					{ p := Point{1, 2}; fmt.Println(p)}					// stack		x-> heap
func f5() []int 			{ s := make([]int, 10); return s }					// heap?
func f6() int 				{ s := make([]int, 10); return len(s) } 			// stack?
func f7(n int) int 			{ s := make([]int, n); return len(s) }  			// heap   x -> stack
func f8() func() int 		{ x := 0; return func() int { x++; return x} }		//heap
func f9(p *Point)			{ global = p }										//stack  x -> heap
func f10() byte 			{ var big [200_000]byte; return big[0] }			// stack? not 10000000 limit. -> heap

func main() {
	p := Point{3,4}
	f9(&p)
}

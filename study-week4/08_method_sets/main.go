// STUDY FILE (theory: method sets + embedding table + method values).
// Different names than the lecture. Run: go run ./08_method_sets
//
// TABLE (what the outer struct's method set contains)
//
//	S embeds T   -> method set of S: methods of T      | of *S: methods of T and *T
//	S embeds *T  -> method set of S: methods of T, *T  | of *S: methods of T and *T
package main

import "fmt"

type Starter interface{ Start() string }

type Engine struct{}

func (Engine) Describe() string { return "an engine" } // value receiver
func (*Engine) Start() string   { return "started" }   // pointer receiver

type Car struct{ Engine }    // embeds Engine by VALUE
type Truck struct{ *Engine } // embeds *Engine (pointer)

var _ Starter = &Car{} // A: compiles (*Car has Engine's *T methods)
// var _ Starter = Car{}                    // B: ERROR, Car's set has only Describe
var _ Starter = Truck{Engine: &Engine{}}  // C: compiles (embedded pointer -> both sets)
var _ Starter = &Truck{Engine: &Engine{}} // D: compiles

// ---- method value: receiver is copied when the method value is created ----
type Rect struct{ W, H float64 }

func (r Rect) Area() float64 { return r.W * r.H }

// ---- Dog vs Cat style: the classic interface/pointer-receiver trap ----
type Speaker interface{ Speak() string }
type Cat struct{}
type Dog struct{}

func (Cat) Speak() string  { return "meow" }
func (*Dog) Speak() string { return "woof" }

func main() {
	var s Speaker
	s = Cat{}  // OK
	s = &Cat{} // OK: *Cat has Cat's methods too
	s = &Dog{} // OK
	// s = Dog{} // ERROR: Dog does not implement Speaker (method Speak has pointer receiver)
	fmt.Println(s.Speak()) // woof

	d := Dog{}
	fmt.Println(d.Speak()) // woof: a direct call on a variable works (Go adds &d)
	// ...but inside an interface the value is an unaddressable copy, so it can't.

	r := Rect{3, 4}
	f := r.Area // receiver r is COPIED right now (W=3)
	r.W = 10
	fmt.Println(f(), r.Area()) // 12 40

	c := Car{}
	fmt.Println(c.Describe(), c.Start()) // an engine started  (c is addressable)
}

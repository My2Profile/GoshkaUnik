// STUDY FILE (analog of Task 3: "compile or not?").
// Run: go run ./03_addressability
//
// THE TWO QUESTIONS to ask for every method call x.M():
//
//  1. Is M a value receiver or a pointer receiver?
//
//  2. Is x ADDRESSABLE (lives in memory Go can point to)?
//
//     value receiver   -> works on anything (it just copies x)
//     pointer receiver -> needs &x, so x must be addressable (or already a pointer)
//
//     Addressable:     variables, slice elements s[i], fields of addressable
//     structs, *p
//     NOT addressable: map elements m[k], function return values, literals T{}
package main

import "fmt"

type Box struct{ n int }

func (b Box) Read() int { return b.n } // value receiver
func (b *Box) Put()     { b.n++ }      // pointer receiver

func NewBox() Box { return Box{} }

func main() {
	var b Box
	b.Put()               // OK  : b is a variable (addressable); Go rewrites to (&b).Put()
	(&b).Put()            // OK  : explicit pointer
	fmt.Println(b.Read()) // 2

	Box{}.Read() // OK  : value receiver works on a non-addressable literal
	// Box{}.Put()      // ERROR: cannot take address of a literal
	// NewBox().Put()   // ERROR: return values are not addressable

	m := map[string]Box{"a": {}}
	m["a"].Read() // OK  : value receiver
	// m["a"].Put()     // ERROR: map elements are not addressable
	// Why? A map may move entries when it grows, so a pointer to one could dangle.

	// Workaround for the map case: copy out, modify, store back.
	tmp := m["a"]
	tmp.Put()
	m["a"] = tmp
	fmt.Println(m["a"].Read()) // 1

	s := []Box{{}}
	s[0].Put()               // OK  : slice elements are addressable
	fmt.Println(s[0].Read()) // 1

	pm := map[string]*Box{"a": {}}
	pm["a"].Put()               // OK  : the map element is a *pointer*; calling through it is fine
	fmt.Println(pm["a"].Read()) // 1
}

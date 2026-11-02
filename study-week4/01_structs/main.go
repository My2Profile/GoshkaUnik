// STUDY FILE (analog of Task 1: struct creation, value vs pointer, copying).
// Domain here is Pet, not Student. Run: go run ./01_structs
//
// BIG PICTURE
//   - A struct is a typed bundle of named fields. No classes, no constructors.
//   - Assigning/passing a struct COPIES every field.
//   - A slice field is copied as a slice HEADER (pointer+len+cap), so the
//     backing array stays shared between copies.
package main

import "fmt"

type Pet struct {
	Name string
	Age  int
	Toys []string
}

// By value: p is a copy of the caller's struct.
// append builds a NEW header (longer len) and stores it in the COPY only.
func addToyByValue(p Pet, toy string) {
	p.Toys = append(p.Toys, toy)
}

// By pointer: we write through the pointer into the caller's struct.
func addToyByPointer(p *Pet, toy string) {
	p.Toys = append(p.Toys, toy)
}

// Even by value, writing to an existing ELEMENT hits the shared array.
func renameFirstToyByValue(p Pet, toy string) {
	p.Toys[0] = toy
}

// A "constructor" in Go is just a normal function, usually named NewX.
func NewPet(name string, age int) *Pet {
	return &Pet{Name: name, Age: age}
}

func main() {
	// ---- 4 ways to create a struct ----
	var p1 Pet                                     // 1. zero value: "" 0 nil
	p2 := Pet{Name: "Rex", Age: 3}                 // 2. named fields (recommended)
	p3 := Pet{"Tom", 2, []string{"ball", "mouse"}} // 3. positional: ALL fields, in order
	p4 := &Pet{Name: "Luna"}                       // 4. pointer to a new struct
	p5 := NewPet("Milo", 1)                        // bonus: constructor function

	fmt.Printf("%+v\n", p1) // {Name: Age:0 Toys:[]}
	fmt.Printf("%+v\n", p2) // {Name:Rex Age:3 Toys:[]}
	fmt.Printf("%+v\n", p3) // {Name:Tom Age:2 Toys:[ball mouse]}
	fmt.Printf("%+v\n", p4) // &{Name:Luna Age:0 Toys:[]}   <- & shows it is a pointer
	fmt.Printf("%+v\n", p5)

	p4.Age = 5 // Go auto-dereferences: same as (*p4).Age = 5
	fmt.Println(p4.Age)

	// ---- by value vs by pointer ----
	fmt.Println("--- append ---")
	addToyByValue(p2, "bone")
	fmt.Println(p2.Toys) // [] : the copy got the new header, we did not
	addToyByPointer(&p2, "bone")
	fmt.Println(p2.Toys) // [bone]

	fmt.Println("--- element write by value ---")
	renameFirstToyByValue(p3, "STICK")
	fmt.Println(p3.Toys) // [STICK mouse] : shared backing array!

	// ---- copying a struct with a slice field ----
	fmt.Println("--- struct copy ---")
	a := Pet{Name: "A", Age: 1, Toys: []string{"x", "y"}}
	b := a // full copy of fields; Toys header copied, array shared
	b.Name = "B"
	b.Age = 99
	b.Toys[0] = "CHANGED"
	fmt.Println(a) // {A 1 [CHANGED y]}  Name/Age untouched, Toys changed
	fmt.Println(b) // {B 99 [CHANGED y]}

	// ---- comparing ----
	// a == b would NOT compile: Pet contains a slice (not comparable).
}

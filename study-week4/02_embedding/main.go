// STUDY FILE (analog of Task 2: embedding, shadowing, name conflict).
// Domain: Device / Phone / Laptop instead of Person / Student / Teacher.
// Run: go run ./02_embedding
//
// BIG PICTURE
//   - Go has no inheritance. "Phone HAS a Device", it is not "a kind of" Device.
//   - Embedding (a field with only a type, no name) PROMOTES the inner
//     type's fields and methods to the outer type as shortcuts.
//   - The embedded field is still a real field named after the type.
package main

import "fmt"

type Device struct {
	Name  string
	Owner string
}

func (d Device) Info() string {
	return d.Name + " (owner: " + d.Owner + ")"
}

// Phone embeds Device and adds a field. Info() is promoted for free.
type Phone struct {
	Device
	Number string
}

// Laptop embeds Device AND declares its own Info(): the outer one wins (shadowing).
type Laptop struct {
	Device
	RAM int
}

func (l Laptop) Info() string {
	return l.Device.Info() + fmt.Sprintf(", RAM %dGB", l.RAM)
}

// ---- name conflict (participation-task analog) ----
type SerialA struct{ ID int }
type SerialB struct{ ID int }

// Declaring Gadget is fine. Only USING the ambiguous name is an error.
type Gadget struct {
	SerialA
	SerialB
}

// ---- embedding is NOT inheritance ----
type Vehicle struct{ Model string }

func (v Vehicle) Sound() string    { return "..." }
func (v Vehicle) Describe() string { return v.Model + " goes " + v.Sound() }

type Car struct{ Vehicle }

func (c Car) Sound() string { return "vroom" }

func main() {
	ph := Phone{Device: Device{Name: "Pixel", Owner: "Aya"}, Number: "+7 700"}
	lp := Laptop{Device: Device{Name: "ThinkPad", Owner: "Bek"}, RAM: 16}

	fmt.Println(ph.Name)        // promoted field: Pixel
	fmt.Println(ph.Device.Name) // full path also works: Pixel
	fmt.Println(ph.Info())      // promoted method: Pixel (owner: Aya)

	fmt.Println(lp.Info())        // Laptop's own (shadows): ThinkPad (owner: Bek), RAM 16GB
	fmt.Println(lp.Device.Info()) // the "hidden" inner one: ThinkPad (owner: Bek)

	// Name conflict:
	g := Gadget{SerialA{1}, SerialB{2}}
	// fmt.Println(g.ID)   // COMPILE ERROR: ambiguous selector g.ID
	fmt.Println(g.SerialA.ID, g.SerialB.ID) // fix: full paths -> 1 2

	// Not inheritance: Describe lives on Vehicle, so inside it v.Sound() is
	// ALWAYS Vehicle.Sound. There is no "virtual dispatch" to Car.Sound.
	c := Car{Vehicle{Model: "Civic"}}
	fmt.Println(c.Sound())         // vroom
	fmt.Println(c.Vehicle.Sound()) // ...
	fmt.Println(c.Describe())      // Civic goes ...   (NOT "vroom")

	// Embedding a POINTER (*Logger) works the same, but if the pointer is nil,
	// calling a promoted method that dereferences it panics.
}

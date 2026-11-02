// STUDY FILE (analog of Task 4: bank account).
// Domain: FuelTank instead of BankAccount. Run: go run ./04_fueltank
//
// BIG PICTURE
//   - Exported (Capital) field = public API. unexported (lowercase) = private
//     to the PACKAGE. Callers must use methods, so the type guards its own rules.
//   - Methods that change state use pointer receivers.
//   - Consistency rule (Go FAQ): if one method needs a pointer receiver, make
//     ALL of them pointer receivers, even read-only ones like Level().
//   - Errors are ordinary values: return error, nil means success.
package main

import (
	"errors"
	"fmt"
)

type FuelTank struct {
	Owner string  // exported
	level float64 // unexported: only code in this package can touch it
}

var ErrNotEnough = errors.New("not enough fuel")

func NewFuelTank(owner string, level float64) *FuelTank {
	return &FuelTank{Owner: owner, level: level}
}

func (t *FuelTank) Fill(amount float64) {
	t.level += amount
}

func (t *FuelTank) Drain(amount float64) error {
	if amount > t.level {
		return fmt.Errorf("drain %.1f from %s: %w", amount, t.Owner, ErrNotEnough)
	}
	t.level -= amount
	return nil
}

// Read-only, but pointer receiver anyway (consistency rule).
func (t *FuelTank) Level() float64 {
	return t.level
}

func main() {
	t := NewFuelTank("truck", 10)
	t.Fill(5)
	fmt.Println(t.Level()) // 15

	if err := t.Drain(100); err != nil {
		fmt.Println("error:", err)                // drain 100.0 from truck: not enough fuel
		fmt.Println(errors.Is(err, ErrNotEnough)) // true
	}
	_ = t.Drain(5)
	fmt.Println(t.Level()) // 10

	// NOTE: here everything is in package main, so t.level would still compile.
	// The "private" protection only bites when FuelTank lives in ANOTHER package.

	// ---- "add 5% to each one in a loop" and the range-copy trap ----
	tanks := []FuelTank{*NewFuelTank("a", 100), *NewFuelTank("b", 200), *NewFuelTank("c", 300)}

	for _, tank := range tanks { // tank is a COPY of each element
		tank.Fill(tank.Level() * 0.05)
	}
	fmt.Println(tanks[0].Level(), tanks[1].Level(), tanks[2].Level()) // 100 200 300 : unchanged!

	for i := range tanks { // index -> real element, addressable
		tanks[i].Fill(tanks[i].Level() * 0.05)
	}
	fmt.Println(tanks[0].Level(), tanks[1].Level(), tanks[2].Level()) // 105 210 315

	// Alternative: store pointers, then even the range variable points at the original.
	ptrs := []*FuelTank{NewFuelTank("x", 100), NewFuelTank("y", 200)}
	for _, p := range ptrs {
		p.Fill(p.Level() * 0.05)
	}
	fmt.Println(ptrs[0].Level(), ptrs[1].Level()) // 105 210
}

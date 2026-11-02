// STUDY FILE (analog of Task 6: Shapes + interface + receivers).
// Domain: Payment methods instead of Shapes. Run: go run ./06_interfaces
//
// BIG PICTURE
//   - An interface is a list of method signatures. A type satisfies it
//     automatically if its METHOD SET contains them all (no "implements").
//   - Method set of T  = value-receiver methods only.
//     Method set of *T = value AND pointer receiver methods.
//   - So a type with pointer-receiver methods can only be stored in an
//     interface as a POINTER (&T{}).
package main

import "fmt"

type Payment interface {
	Name() string
	Fee(amount float64) float64
}

// Card: value receivers -> both Card{} and &Card{} satisfy Payment.
type Card struct{ Balance float64 }

func (c Card) Name() string               { return "card" }
func (c Card) Fee(amount float64) float64 { return amount * 0.01 }

// TopUp must really change the card, so it needs a pointer receiver.
// Mixing is allowed because TopUp is not part of the Payment interface.
func (c *Card) TopUp(x float64) { c.Balance += x }

// Crypto: pointer receivers -> only *Crypto satisfies Payment.
type Crypto struct{ Network string }

func (c *Crypto) Name() string               { return "crypto/" + c.Network }
func (c *Crypto) Fee(amount float64) float64 { return 2.5 }

// Compile-time proof that *Crypto implements Payment (costs nothing at runtime).
var _ Payment = (*Crypto)(nil)

func main() {
	payments := []Payment{
		Card{Balance: 10},       // OK: value in method set
		&Crypto{Network: "eth"}, // must be a pointer
	}
	// payments = append(payments, Crypto{Network: "btc"})
	// ^ COMPILE ERROR: Crypto does not implement Payment (method Name has pointer receiver)
	// Why: the value stored inside an interface is a copy that is not addressable,
	// so Go cannot take its address to call a pointer-receiver method.

	for _, p := range payments {
		fmt.Printf("%-10s fee on 200 = %.2f\n", p.Name(), p.Fee(200))
	}

	card := Card{Balance: 10}
	card.TopUp(5)             // card is a variable -> addressable -> Go uses (&card).TopUp
	fmt.Println(card.Balance) // 15
}

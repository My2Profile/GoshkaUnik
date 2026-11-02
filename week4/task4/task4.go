package main

import (
	"errors"
	"fmt"
)

type Bank struct {
	Owner   string
	balance float64
}

func (b *Bank) Deposit(amount float64) {
	b.balance += amount
}

func (b *Bank) Withdraw(amount float64) error {
	if b.balance < amount {
		return errors.New("Insufficient funds")
	}
	b.balance -= amount
	return nil
}

func (b Bank) Balance() float64 {
	return b.balance
}

func main() {

	b1 := Bank{Owner: "Alice", balance: 1000.0}

	fmt.Println(b1)

	for i := 0; i < 5; i++ {
		b1.Deposit(b1.Balance() * 0.05)
	}

	fmt.Println("Balance after 5 months of 5% interest:", b1.Balance())

	err := b1.Withdraw(2000.0)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

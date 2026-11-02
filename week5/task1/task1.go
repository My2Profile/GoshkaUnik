package main

import (
	"errors"
	"fmt"
)

// 1. Declare Notifier interface
type Notifier interface {
	Notify(msg string) error
}

// 4. Compile-time interface satisfaction checks
var _ Notifier = Email{}
var _ Notifier = SMS{}
var _ Notifier = &Telegram{}

// 2. Implementations
type Email struct {
	Address string
}

func (e Email) Notify(msg string) error {
	fmt.Printf("Sending Email to %s: %s\n", e.Address, msg)
	return nil
}

type SMS struct {
	Phone string
}

func (s SMS) Notify(msg string) error {
	if len(msg) > 20 {
		return errors.New("SMS error: message too long (>20 chars)")
	}
	fmt.Printf("Sending SMS to %s: %s\n", s.Phone, msg)
	return nil
}

type Telegram struct {
	Chat      string
	SentCount int
}

// Must use pointer receiver to mutate SentCount across calls
func (t *Telegram) Notify(msg string) error {
	t.SentCount++
	fmt.Printf("Sending Telegram to %s [#%d]: %s\n", t.Chat, t.SentCount, msg)
	return nil
}

// 3. Broadcast function
func Broadcast(ns []Notifier, msg string) []error {
	var errs []error
	for _, n := range ns {
		if err := n.Notify(msg); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func main() {
	notifiers := []Notifier{
		Email{Address: "user@example.com"},
		SMS{Phone: "+123456789"},
		&Telegram{Chat: "@golang_channel"},
	}

	errs := Broadcast(notifiers, "Hello team!")
	if len(errs) > 0 {
		fmt.Println("Broadcast encountered errors:", errs)
	}
}

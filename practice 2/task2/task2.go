package main

import "fmt"

func main() {

	scores := map[string]int{
		"Alice": 100,
		"Bob":   70,
	}

	if value, ok := scores["Charlie"]; ok {
		fmt.Printf("Charlie is in the list: %d", value)
	} else {
		fmt.Println("Charlie is missing")
	}

	if value, ok := scores["Bob"]; ok {
		fmt.Printf("Bob is in the list: %d", value)
	} else {
		fmt.Println("Bob is missing")
	}

	manageInventory()

}

func manageInventory() {

	inventory := map[string]int{
		"apples":  10,
		"bananas": 5,
	}

	fmt.Println("\n", "Current apples stock:", inventory["apples"])

	fmt.Println("Current Inventory: ", inventory)

	inventory["bananas"] = 12

	inventory["oranges"] = 8

	delete(inventory, "apples")

	fmt.Println("Updated inventory:", inventory)
}

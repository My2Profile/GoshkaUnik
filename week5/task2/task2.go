package main

import (
	"fmt"
)

func describe(v any) string {
	switch x := v.(type) {
	case nil:
		return "nil"
	case int, int64:
		return fmt.Sprintf("integer %v", x)
	case string:
		return fmt.Sprintf("string of %d bytes", len(x)) // x is a string here
	case []int:
		return fmt.Sprintf("slice with %d ints", len(x)) // x is a []int here
	case error:
		return "error: " + x.Error()
	case fmt.Stringer:
		return "Stringer: " + x.String()
	case bool:
		return fmt.Sprintf("boolean %v", x)
	case map[string]int:
		return fmt.Sprintf("map with %d entries", len(x))
	default:
		return fmt.Sprintf("unknown type %T", x)
	}
}

type Celsius float64

func (c Celsius) String() string { return fmt.Sprintf("%.1f°C", float64(c)) }

func main() {

	fmt.Println(describe(true))  // Output: boolean true
	fmt.Println(describe(false)) // Output: boolean false

	fmt.Println(describe(nil)) // Output: nil

	fmt.Println(describe(7))                       // Output: integer 7
	fmt.Println(describe(int64(7)))                // Output: integer 7
	fmt.Println(describe("привет"))                // Output: string of 12 bytes
	fmt.Println(describe([]int{1, 2, 3}))          // Output: slice with 3 ints
	fmt.Println(describe(fmt.Errorf("disk full"))) // Output: error: disk full
	fmt.Println(describe(Celsius(21.5)))           // Output: Stringer: 21.5°C
	fmt.Println(describe(3.14))                    // Output: unknown type float64
	fmt.Println(describe(struct{}{}))

	var t any = Celsius(36.6)
	if st, ok := t.(fmt.Stringer); ok {
		fmt.Println("Stringer:", st.String()) // Stringer: 36.6°C
	}
	// Output: unknown type struct {}

}

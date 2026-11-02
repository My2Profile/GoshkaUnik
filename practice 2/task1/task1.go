package main

import (
	"bytes"
	"fmt"
	"strings"
	"unicode/utf8"
)

func main() {
	s := "hêllo"
	for i := range s {
		fmt.Printf("position %d: %c\n", i, s[i])
	}

	for i, r := range s {
		fmt.Printf("position: %d: %c\n", i, r)
	}

	trim_1 := strings.TrimRight("123oxo", "xo")
	fmt.Println("Original: 123oxo", "\n", "TrimRight:", trim_1)

	trim_2 := strings.TrimLeft("xoox123oxox", "ox")
	fmt.Println("Original: xoox123oxox", "\n", "TrimLeft:", trim_2)

	trim_3 := strings.TrimPrefix("xoo123oxox", "xoo")
	fmt.Println("Original: xoo123oxox", trim_3, "\n", "TrimPrefix:", trim_3)

	trim_4 := strings.TrimSuffix("123oxox", "ox")
	fmt.Println("Original: 123oxox", trim_4, "\n", "TrimSuffix:", trim_4)

	words := []string{"lol", "123"}

	fmt.Println(concat(words))

	strings_with_trails := "String  "
	fmt.Printf("processed String: %q", processPayloadOptimized([]byte(strings_with_trails)))

	fmt.Println(utf8.RuneCountInString(s))

}

func concat(values []string) string {

	total_length := 0

	for _, v := range values {
		total_length += len(v)
	}

	var sb strings.Builder
	sb.Grow(total_length)

	for _, j := range values {
		sb.WriteString(j)
	}
	return sb.String()

}

func processPayload(data []byte) string {
	// Unnecessary type conversions
	trimmed := string(bytes.TrimSpace(data))
	clean := strings.ReplaceAll(trimmed, "\r", "")
	return clean
}

func processPayloadOptimized(data []byte) string {
	trimmed := bytes.TrimSpace(data)
	clean := string(bytes.ReplaceAll(trimmed, []byte("\r"), []byte("")))
	return clean
}

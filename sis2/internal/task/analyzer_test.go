package task

import "testing"

func TestTextAnalyzer(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Result
	}{
		{"simple", "hello world", Result{Words: 2, Runes: 11, Lines: 1, UniqueWords: 2}},
		{"case insensitive unique", "Go go GO", Result{Words: 3, Runes: 8, Lines: 1, UniqueWords: 1}},
		{"multi line", "a b\nc d\ne", Result{Words: 5, Runes: 9, Lines: 3, UniqueWords: 5}},
		{"cyrillic counts runes not bytes", "привет мир", Result{Words: 2, Runes: 10, Lines: 1, UniqueWords: 2}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := TextAnalyzer{}.Process(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if *got != tc.want {
				t.Errorf("got %+v, want %+v", *got, tc.want)
			}
		})
	}
}

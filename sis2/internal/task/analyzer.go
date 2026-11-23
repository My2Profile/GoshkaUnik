package task

import (
	"strings"
	"unicode/utf8"
)

type Processor interface {
	Process(input string) (*Result, error)
}

type TextAnalyzer struct{}

func (TextAnalyzer) Process(input string) (*Result, error) {
	words := strings.Fields(input)
	runes := utf8.RuneCountInString(input)
	lines := strings.Count(input, "\n") + 1
	uniqueWords := make(map[string]struct{})
	for _, word := range words {
		uniqueWords[strings.ToLower(word)] = struct{}{}
	}
	return &Result{
		Words:       len(words),
		Runes:       runes,
		Lines:       lines,
		UniqueWords: len(uniqueWords),
	}, nil

}

var _ Processor = TextAnalyzer{}

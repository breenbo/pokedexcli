package internal

import "strings"

func CleanInput(text string) []string {
	words := []string{}

	cleaned_text := strings.ToLower(strings.Trim(text, " "))

	input := strings.SplitSeq(cleaned_text, " ")
	for word := range input {
		words = append(words, word)
	}

	return words
}

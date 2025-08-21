package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		text := scanner.Text()
		cleaned := cleanInput(text)
		fmt.Printf("Your command was: %s\n", cleaned[0])
	}

}

func cleanInput(text string) []string {
	words := []string{}

	cleaned_text := strings.ToLower(strings.Trim(text, " "))

	input := strings.SplitSeq(cleaned_text, " ")
	for word := range input {
		words = append(words, word)
	}

	return words
}

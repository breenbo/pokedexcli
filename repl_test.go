package main

import (
	"github.com/breenbo/pokedexcli/internal"
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "hello world",
			expected: []string{"hello", "world"},
		},
		{
			input:    "    hello world   ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "HELLO WORLD",
			expected: []string{"hello", "world"},
		},
	}

	for _, c := range cases {
		actual := internal.CleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf("not same count of words")
		}

		for i := range actual {
			word := actual[i]
			expected := c.expected[i]

			if word != expected {
				t.Errorf("word are not matching: %s - %s", word, expected)
			}
		}
	}
}

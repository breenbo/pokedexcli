package internal

import "fmt"

func PrintPokemon(pokemon Pokemon) {
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %v\n", pokemon.Height)
	fmt.Printf("Weight: %v\n", pokemon.Weight)
	fmt.Print("Stats:\n")
	for _, value := range pokemon.Stats {
		for k, v := range value {
			fmt.Printf(" - %v: %v\n", k, v)
		}
	}
	fmt.Print("Types: \n")
	for _, t := range pokemon.Types {
		fmt.Printf(" - %s\n", t)
	}
}

func PrintResults(results []Result) {
	// print the locations from slice
	locationMsg := "\n"
	for _, loc := range results {
		locationMsg += fmt.Sprintf(" - %s\n", loc.Name)
	}
	locationMsg += "\n"

	fmt.Print(locationMsg)
}

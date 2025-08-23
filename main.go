package main

import (
	"bufio"
	"fmt"
	"github.com/breenbo/pokedexcli/internal"
	"math/rand"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(config *Config, argument string) error
}

var commands map[string]cliCommand

type Config struct {
	Next     string
	Previous string
	BaseUrl  string
}

var config Config
var pokedex map[string]internal.Pokemon

func main() {
	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Help for Pokedex",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Display 20 pokemon location",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Display 20 pokemon location back",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Search for pokemon inside a map area - need a valid area name",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "try to catch a pokemon - need a valid pokemon name",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "inspect a pokemon in your pokedex - need a valid pokemon name",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "inspect your pokedex - need a valid pokemon name",
			callback:    commandPokedex,
		},
	}

	locationURL := "https://pokeapi.co/api/v2/location-area"
	config = Config{
		BaseUrl:  "https://pokeapi.co/api/v2",
		Next:     locationURL,
		Previous: "",
	}
	pokedex = make(map[string]internal.Pokemon)

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		text := scanner.Text()
		cleaned := internal.CleanInput(text)
		command := cleaned[0]

		argument := ""
		if len(cleaned) > 1 {
			argument = cleaned[1]
		}

		if value, ok := commands[command]; ok {
			err := value.callback(&config, argument)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Print("Unknown command\n")
		}
	}

}

func commandExit(config *Config, argument string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)

	return nil
}

func commandHelp(config *Config, argument string) error {
	helpMsg := "Welcome to the Pokedex!\n\nUsage:\n\n"

	for _, command := range commands {
		helpMsg += fmt.Sprintf("%s: %s\n", command.name, command.description)
	}

	helpMsg += "\n"

	fmt.Print(helpMsg)

	return nil

}

func commandMap(config *Config, argument string) error {
	res, err := internal.FetchLocation(config.Next)

	if err != nil {
		return err
	}
	// set up the config
	config.Next = res.Next
	config.Previous = res.Previous

	internal.PrintResults(res.Results)

	return nil
}

func commandMapb(config *Config, argument string) error {
	if config.Previous == "" {
		fmt.Print("you're on the first page\n")
		return nil
	}

	res, err := internal.FetchLocation(config.Previous)

	if err != nil {
		return err
	}
	// set up the config
	config.Next = res.Next
	config.Previous = res.Previous

	internal.PrintResults(res.Results)

	return nil
}

func commandExplore(config *Config, areaName string) error {
	if areaName == "" {
		fmt.Print("you must add a valid area name\n")
		return nil
	}

	fmt.Printf("Exploring %s...\n", areaName)

	fullUrl := config.BaseUrl + "/location-area"

	res, err := internal.FetchExplore(fullUrl, areaName)
	if err != nil {
		return err
	}

	if len(res) > 0 {
		fmt.Print("Found Pokemon:")
		internal.PrintResults(res)
	} else {
		fmt.Print("No pokemon found\n")
	}

	return nil
}

func commandCatch(config *Config, pokemonName string) error {
	if pokemonName == "" {
		fmt.Print("you must add valid pokemon name\n")
		return nil
	}

	if _, ok := pokedex[pokemonName]; ok {
		fmt.Printf("%s already in your pokedex\n", pokemonName)
		return nil
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
	fullUrl := config.BaseUrl + "/pokemon"

	pokemon, err := internal.FetchPokemon(fullUrl, pokemonName)
	if err != nil {
		return err
	}

	difficulty := pokemon.Base_experience
	randRoll := rand.Intn(difficulty + 100)

	if randRoll > difficulty {
		fmt.Printf("%s escaped!\n", pokemonName)
	} else {
		fmt.Printf("%s was caught!\n", pokemonName)
		pokedex[pokemonName] = pokemon
	}

	return nil
}

func commandInspect(config *Config, pokemonName string) error {
	if pokemonName == "" {
		fmt.Print("you must add valid pokemon name\n")
		return nil
	}

	if _, ok := pokedex[pokemonName]; !ok {
		fmt.Printf("%s not in your pokedex\n", pokemonName)
		return nil
	}

	fullUrl := config.BaseUrl + "/pokemon"

	pokemon, err := internal.FetchPokemon(fullUrl, pokemonName)
	if err != nil {
		return err
	}

	internal.PrintPokemon(pokemon)

	return nil
}

func commandPokedex(config *Config, arg string) error {
	if len(pokedex) == 0 {
		fmt.Print("No pokemon in your pokedex\n")
		return nil
	}

	fmt.Print("Your Pokedex:\n")
	for _, p := range pokedex {
		fmt.Printf(" - %v\n", p.Name)
	}

	return nil
}

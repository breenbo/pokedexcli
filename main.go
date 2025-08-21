package main

import (
	"bufio"
	"fmt"
	"github.com/breenbo/pokedexcli/internal"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(config *Config) error
}

var commands map[string]cliCommand

type Config struct {
	Next     string
	Previous string
}

var config Config

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
	}

	locationURL := "https://pokeapi.co/api/v2/location-area"
	config = Config{
		Next:     locationURL,
		Previous: "",
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		text := scanner.Text()
		cleaned := internal.CleanInput(text)
		command := cleaned[0]

		if value, ok := commands[command]; ok {
			err := value.callback(&config)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Print("Unknown command\n")
		}
	}

}

func commandExit(config *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)

	return nil
}

func commandHelp(config *Config) error {
	helpMsg := "Welcome to the Pokedex!\n\nUsage:\n\n"

	// TODO: sort func by alpha order ?

	for _, command := range commands {
		helpMsg += fmt.Sprintf("%s: %s\n", command.name, command.description)
	}

	helpMsg += "\n"

	fmt.Print(helpMsg)

	return nil

}

func commandMap(config *Config) error {
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

func commandMapb(config *Config) error {
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

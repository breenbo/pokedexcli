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
	callback    func(config *Config, argument string) error
}

var commands map[string]cliCommand

type Config struct {
	Next     string
	Previous string
	BaseUrl  string
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
		"explore": {
			name:        "explore",
			description: "Search for pokemon inside a map area - need a valid area name",
			callback:    commandExplore,
		},
	}

	locationURL := "https://pokeapi.co/api/v2/location-area"
	config = Config{
		BaseUrl:  locationURL,
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

	// TODO: sort func by alpha order ?

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

	res, err := internal.FetchExplore(config.BaseUrl, areaName)
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

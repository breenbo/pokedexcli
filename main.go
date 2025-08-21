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
	callback    func() error
}

var commands map[string]cliCommand

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

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		text := scanner.Text()
		cleaned := internal.CleanInput(text)
		command := cleaned[0]

		if value, ok := commands[command]; ok {
			err := value.callback()
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Print("Unknown command\n")
		}
	}

}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)

	return nil
}

func commandHelp() error {
	helpMsg := "Welcome to the Pokedex!\n\nUsage:\n\n"

	// TODO: sort func by alpha order ?

	for _, command := range commands {
		helpMsg += fmt.Sprintf("%s: %s\n", command.name, command.description)
	}

	helpMsg += "\n"

	fmt.Print(helpMsg)

	return nil

}

func commandMap() error {
	res, err := internal.FetchLocation()

	if err != nil {
		return err
	}

	locationMsg := ""
	for _, loc := range res.Results {
		locationMsg += fmt.Sprintf("%s\n", loc.Name)
	}
	locationMsg += "\n"

	fmt.Print(locationMsg)

	return nil
}

func commandMapb() error {
	return nil
}

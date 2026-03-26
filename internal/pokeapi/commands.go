package pokeapi

import (
	"fmt"
	"os"
)

type CliCommand struct {
	name        string
	description string
	callback    func(cfg *Config, args ...string) error
}

type Config struct {
	client 		Client
	Next 		*string
	Previous 	*string
}

func InitCommands() {
	commands = map[string]CliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name: 		"help",
			description: "Display help message",
			callback: commandHelp,
		},
		"map": {
			name:		"map",
			description: "Display 20 location names",
			callback: commandMap,
		},
		"mapb": {
			name:		"mapb",
			description: "Display previous 20 location names",
			callback: commandMapb,
		},
		"explore": {
			name: 		"explore",
			description: "Display pokemon in a specific location area",
			callback: commandExplore,
		},
	}
}

func commandExit(cfg *Config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *Config, args ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for key, value := range commands {
		fmt.Printf("%s: %s\n", key, value.description)
	}
	return nil		
}

func commandMap(cfg *Config, args ...string) error {
	if cfg.Next != nil{
		return mapGetHelper(cfg, *cfg.Next)
	} else {
		return mapGetHelper(cfg, "https://pokeapi.co/api/v2/location-area")
	}
}

func commandMapb(cfg *Config, args ...string) error {
	if cfg.Previous != nil{
		return mapGetHelper(cfg, *cfg.Previous)
	} else {
		fmt.Println("you're on the first page")
		return nil
	}
}

func mapGetHelper(cfg *Config, url string) error {
	locationAreas, err := cfg.client.LocationAreaGET(url)
	fmt.Println("debug: LocationAreaGET returned, count:", locationAreas.Count)
    fmt.Println("debug: results length:", len(locationAreas.Results))
	if err != nil {
		return err
	}
	cfg.Next = locationAreas.Next
	cfg.Previous = locationAreas.Previous
	for _, result := range locationAreas.Results {
		fmt.Println(result.Name)
	}
	fmt.Println("debug: mapGetHelper returning")
	return nil
}

func commandExplore(cfg *Config, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("Expecting 1 argument")
	}
	location, err := cfg.client.LocationGET(args[0])
	if err != nil {
		return err
	}
	for _, encounter := range location.PokemonEncounters {
		fmt.Println(encounter.Pokemon.Name)
	}
	
	return nil
}
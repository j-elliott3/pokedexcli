package main

import (
	"strings"
	"bufio"
	"os"
	"fmt"
	"time"
	"github.com/j-elliott3/pokedexcli/internal/pokeapi"
)

var commands map[string]pokeapi.CliCommand

func CleanInput(text string) []string {
	words := strings.Fields(strings.ToLower(text))
	return words
}

func StartREPL() {
	InitCommands()
	interval := 5 * time.Minute
	client := pokeapie.NewClient(interval)
	cfg := &pokeapi.Config{
		client: client
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()

		cleanText := CleanInput(scanner.Text())
		if len(cleanText) == 0 {
		continue
		}
		command, ok := commands[cleanText[0]]
		if ok {
			err := command.callback(cfg)
			if err != nil {
				fmt.Println(err)
				continue
			}
		} else {
			fmt.Println("Unknown command")
			continue
		}
	}
}
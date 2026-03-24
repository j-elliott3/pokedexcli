package main

import (
	"strings"
	"bufio"
	"os"
	"fmt"
	"github.com/j-elliott3/pokedexcli/internal/pokecache"
)

var commands map[string]cliCommand

func CleanInput(text string) []string {
	words := strings.Fields(strings.ToLower(text))
	return words
}

func StartREPL() {
	initCommands()
	cfg := &Config{}
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
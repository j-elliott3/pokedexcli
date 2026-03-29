package pokeapi

import (
	"strings"
	"bufio"
	"os"
	"fmt"
	"time"
)

var commands map[string]CliCommand

func CleanInput(text string) []string {
	words := strings.Fields(strings.ToLower(text))
	return words
}

func StartREPL() {
	InitCommands()
	interval := 5 * time.Minute
	client := NewClient(interval)
	cfg := &Config{
		client: client,
		Pokedex: make(map[string]Pokemon),
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
			err := command.callback(cfg, cleanText[1:]...)
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
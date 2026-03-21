package main

import (
	"strings"
	"bufio"
	"os"
	"fmt"
)

func CleanInput(text string) []string {
	words := strings.Fields(strings.ToLower(text))
	return words
}

func StartREPL() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		cleanText := CleanInput(scanner.Text())
		fmt.Printf("Your command was: %s\n", cleanText[0])
	}
}
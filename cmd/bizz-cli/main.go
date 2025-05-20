package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

func main() {
	fmt.Println("BizzGenie CLI - Interactive Shell")
	fmt.Println("Type 'help' for available commands or 'exit' to quit")

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		input := scanner.Text()
		args := strings.Fields(input)

		if len(args) == 0 {
			continue
		}
		switch strings.ToLower(args[0]) {
		case "exit", "quit":
			fmt.Println("GoodBye")
			return
		case "help":
			fmt.Println("Helpeig")

		default:
			fmt.Println("Default is running")

		}

		printPrompt()

	}

}

func printPrompt() {
	color.New(color.FgHiCyan).Print("bizzcli> ")
}

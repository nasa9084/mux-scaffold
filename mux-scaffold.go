package main

import (
	"fmt"
	"os"

	"github.com/nasa9084/mux-scaffold/command"
)

func main() { os.Exit(exec()) }

func exec() int {
	if len(os.Args) < 2 {
		printHelp()
		return 1
	}

	commandName := os.Args[1]
	commandArgs := []string{}
	if len(os.Args) > 2 {
		commandArgs = os.Args[2:]
	}

	command, ok := command.List[commandName]
	if !ok {
		printHelp()
		return 1
	}

	return command.Exec(commandArgs)
}

func printHelp() {
	fmt.Printf("mux Scaffold\n")
	for name, cmd := range command.List {
		fmt.Printf("\nCommand \"%s\":\n", name)
		fmt.Printf("        %s\n", cmd.Description())
	}
}

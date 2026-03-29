package main

import (
	"fmt"
	"os"
)

func main() {

	appState := &State{
		Config: &struct {
			CurrentUserName string
		}{
			CurrentUserName: "Mike",
		},
	}
	cmdRegistry := NewCommands()
	cmdRegistry.Register("login", HandlerLogin)

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Error: missing command name\n")
		fmt.Fprintf(os.Stderr, "Usage: %s <command> [args]\n", os.Args[0])
		os.Exit(1)
	}
	commandName := os.Args[1]
	commandArgs := os.Args[2:]
	cmdInstance := Command{
		Name: commandName,
		Args: commandArgs,
	}
	if err := cmdRegistry.Run(appState, cmdInstance); err != nil {
		fmt.Printf("Command failed: %V\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nFinal state after command execution: %+v\n", appState.Config)
}

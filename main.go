package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/mpdev25/pokedexcli/blog_aggregator/internal/config"
)

func main() {

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading file: %v\n", err)
	}
	appState := &State{
		Config: &cfg,
	}
	cmdRegistry := NewCommands()
	cmdRegistry.Register("login", HandlerLogin)
	cmdRegistry.Register("setdb", HandlerSetDB)

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

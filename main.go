package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"github.com/mpdev25/pokedexcli/blog_aggregator/internal/config"
	"github.com/mpdev25/pokedexcli/blog_aggregator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func registerHandler(s *state, cmd string, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: %s <name>", cmd)
	}
	name := args[0]
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	})
	if err != nil {
		return fmt.Errorf("could not create user %s: %w", name, err)
	}
	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("could not set user: %w", err)
	}
	fmt.Printf("User '%s' created successfully\n", name)
	log.Printf("User Data: %+v\n", user)
	return nil
}

func main() {

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading file: %v\n", err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}
	defer db.Close()

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

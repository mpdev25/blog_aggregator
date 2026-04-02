package main

import (
	"fmt"
	"log"

	"github.com/mpdev25/pokedexcli/blog_aggregator/internal/config"
)

type State struct {
	Config *config.Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Handlers map[string]func(*State, Command) error
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("login failed: expected a single argument (username), but got %d", len(cmd.Args))
	}
	username := cmd.Args[0]

	if err := s.Config.SetUser(username); err != nil {
		log.Printf("Failed to set user to %s: %v", username, err)
		return fmt.Errorf("user configuration failed: %w", err)
	}

	fmt.Println("Username has been set to:", s.Config.CurrentUserName)
	return nil
}

func HandlerSetDB(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <db_url>", cmd.Name)
	}
	dbURL := cmd.Args[0]
	if err := s.Config.SetDatabaseURL(dbURL); err != nil {
		return fmt.Errorf("failed to set database URL: %w", err)
	}
	fmt.Printf("Database URL has been updated to: %s\n", dbURL)
	return nil
}

func (c *Commands) Run(s *State, cmd Command) error {
	if c.Handlers == nil {
		c.Handlers = make(map[string]func(*State, Command) error)
	}
	if handler, exists := c.Handlers[cmd.Name]; exists {
		return handler(s, cmd)
	} else {
		return fmt.Errorf("command not found: %s", cmd.Name)
	}

}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	if c.Handlers == nil {
		c.Handlers = make(map[string]func(*State, Command) error)
	}
	c.Handlers[name] = f
}

func NewCommands() *Commands {
	return &Commands{
		Handlers: make(map[string]func(*State, Command) error),
	}
}

package main

import (
	"context"
	"fmt"
	"log"
	"os"
)

//type State struct {
//	Config *config.Config
//}

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

	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		fmt.Printf("User '%s' does not exist in the database.\n", username)
		os.Exit(1)
	}

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

func Reset(s *State, cmd Command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to reset table: %w\n", err)

	}
	fmt.Println("user table reset")
	return nil
}

func users(s *State, cmd Command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to retrieve users: %w\n", err)
	}
	for _, user := range users {
		if user == s.Config.CurrentUserName {
			user = fmt.Sprintf("%s (current)", user)
		}
		fmt.Printf("* %s\n", user)
	}
	return nil
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	if c.Handlers == nil {
		c.Handlers = make(map[string]func(*State, Command) error)
	}
	c.Handlers[name] = f
}

func Agg(s *State, cmd Command) error {
	feedURL := "https://www.wagslane.dev/index.xml"
	for _, URL := range os.Args[1] {
		if err != nil {
			fmt.Errorf("unable to retrieve url %v\n", err)
		}
	}
	if len(os.Args) > 1 {
		fmt.Println("", os.Args[1])
	}
	return nil
}

func NewCommands() *Commands {
	return &Commands{
		Handlers: make(map[string]func(*State, Command) error),
	}
}

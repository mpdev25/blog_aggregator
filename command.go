package main

import (
	"fmt"

	"github.com/mpdev25/pokedexcli/blog_aggregator/internal/config"
)

type State struct {
	Config *config.Config
}

type Command struct {
	Name string
	Args []string
}

func handlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("login failed: expected a single argument (username), but got %d", len(cmd.Args))
	}
	username := cmd.Args[0]
	if s.Config == nil {
		s.Config = &config.Config{}
	}
	s.Config.CurrentUserName = username
	fmt.Println("Username has been set to:", s.Config.CurrentUserName)
	return nil
}

type Commands struct {
	Handlers map[string]func(*State, Command) error
}

func (c *Commands) run(s *State, cmd Command) error {
	if c.Handlers == nil {
		c.Handlers = make(map[string]func(*State, Command) error)
	}
	if handler, exists := c.Handlers[cmd.Name]; exists {
		return handler(s, cmd)
	} else {
		return fmt.Errorf("command not found: %s", cmd.Name)
	}

}

func (c *Commands) register(name string, f func(*State, Command) error) {
	if c.Handlers == nil {
		c.Handlers = make(map[string]func(*State, Command) error)
	}
	c.Handlers[name] = f
}

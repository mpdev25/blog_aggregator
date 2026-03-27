package command

import (
	"fmt"
)

type State struct {
	Pointer *config
}

type Command struct {
	Name        string
	StringSlice []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("login failed: expected a single argument (username), but got %d", len(args))
	}
	username = cmd.args[0]
	s.config.Username = username
	fmy.Println("username has been set")
	return nil
}

type Commands struct {
	Handlers map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {

}

func (c *commands) register(name string, f func(*state, command) error) {

}

package main

import "fmt"

func middlewareLoggedIn(handler func(s *State, cmd Command) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		if s.Config.CurrentUserName == "" {
			return fmt.Errorf("this command requires you to be logged in")
		}
		return handler(s, cmd)
	}
}

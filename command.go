package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mpdev25/pokedexcli/blog_aggregator/internal/database"
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
	resp, err := http.Get("https://www.wagslane.dev/index.xml")

	if err != nil {
		return fmt.Errorf("unable to retrieve url %v\n", err)
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))

	return nil
}

func addfeed(s *State, cmd Command) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("usage: %s addfeed <name> <url>", cmd.Name)
	}
	if s.Config.CurrentUserName == "" {
		return fmt.Errorf("you must be logged in to add a feed")
	}
	user, err := s.db.GetUser(context.Background(), s.Config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("could not find current user: %w", err)
	}
	feed, err := s.db.CreateFeeds(context.Background(), database.CreateFeedsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	})

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "unique violation") {
			return fmt.Errorf("feed with name '%s' or URL '%s' already exists", cmd.Args[0], cmd.Args[1])
		}
		return fmt.Errorf("failed to add feed: %w", err)
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("feed added, but failed to follow automatically: %w", err)
	}

	fmt.Printf("Feed '%s' added successfully\n", feed.Name)
	return nil
}

func feeds(s *State, cmd Command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("could not get feeds: %w", err)
	}

	for _, feed := range feeds {
		fmt.Printf("* %s\n", feed.Name)
		fmt.Printf(" - URL: %s\n", feed.Url)
		fmt.Printf(" - Created by: %s\n", feed.UserName)
	}
	return nil
}

func follow(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}
	url := cmd.Args[0]
	if s.Config.CurrentUserName == "" {
		return fmt.Errorf("you must be logged in to follow a feed")
	}
	user, err := s.db.GetUser(context.Background(), s.Config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("could not get user: %w", err)
	}
	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("could not find feed with URL %s: %w", url, err)
	}
	follow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("could not follow feed: %w", err)
	}
	fmt.Printf("User '%s' is now following feed '%s'\n", follow.UserName, follow.FeedName)
	return nil

}
func following(s *State, cmd Command) error {
	if s.Config.CurrentUserName == "" {
		return fmt.Errorf("you must be logged in to follow a feed")
	}

	follows, err := s.db.GetFeedFollowsForUser(context.Background(), s.Config.CurrentUserName) //user.ID.String())
	if err != nil {
		return fmt.Errorf("could not retrieve feeds: %w", err)
	}
	fmt.Printf("Feed follows for user '%s':\n", s.Config.CurrentUserName)
	for _, follow := range follows {
		fmt.Printf("* %s\n", follow.FeedName)
	}
	return nil
}

func unfollow(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}
	url := cmd.Args[0]
	if s.Config.CurrentUserName == "" {
		return fmt.Errorf("you must be logged in to unfollow a feed")
	}
	user, err := s.db.GetUser(context.Background(), s.Config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("could not get user: %w", err)
	}
	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("could not find feed with URL %s: %w", url, err)
	}

	err = s.db.Delete_Feed_Follow(context.Background(), database.Delete_Feed_FollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("could not unfollow feed: %w", err)
	}
	fmt.Printf("User '%s' is noo longer following feed '%s'\n", user.Name, feed.Url)
	return nil
}

func NewCommands() *Commands {
	return &Commands{
		Handlers: make(map[string]func(*State, Command) error),
	}
}

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mpdev25/pokedexcli/blog_aggregator/internal/database"
)

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

	if len(cmd.Args) < 1 || len(cmd.Args) > 2 {
		return fmt.Errorf("usage: %s <time_between_reqs>", cmd.Name)
	}
	time_between_reqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %v", err)
	}
	fmt.Printf("Collecting feeds every %s ...\n", time_between_reqs)
	ticker := time.NewTicker(time_between_reqs)

	for ; ; <-ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			log.Printf("Scraper error: %v", err)
		}
	}

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

	follows, err := s.db.GetFeedFollowsForUser(context.Background(), s.Config.CurrentUserName)
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

func scrapeFeeds(s *State) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("could not find next feed: %w", err)
	}
	log.Printf("Found a feed to fetch: %s", nextFeed.Url)
	return ScrapeFeed(s.db, database.Feed{
		ID:     nextFeed.ID,
		Url:    nextFeed.Url,
		Name:   nextFeed.Name,
		UserID: nextFeed.UserID,
	})

}
func ScrapeFeed(db *database.Queries, feed database.Feed) error {
	feedData, err := fetchFeed(context.Background(), feed.Url)

	if err != nil {
		return fmt.Errorf("could not parse feed %s: %w", feed.Url, err)
	}

	_, err = db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return fmt.Errorf("could not mark feed as fetched: %w", err)
	}

	for _, item := range feedData.Channel.Item {
		description := item.Description

		parsedTime, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			parsedTime, err = time.Parse(time.RFC1123, item.PubDate)
		}
		publishedAt := sql.NullTime{}
		if err == nil {
			publishedAt = sql.NullTime{
				Time:  parsedTime,
				Valid: true,
			}
		}
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),

			Title: sql.NullString{
				String: item.Title,
				Valid:  item.Title != "",
			},
			Url: item.Link,

			Description: sql.NullString{
				String: description,
				Valid:  description != "",
			},
			PublishedAt: publishedAt,

			FeedID: feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Database error saving post: %v", err)
		}

	}

	return nil
}

func browse(s *State, cmd Command) error {
	limit := 2
	if len(cmd.Args) > 0 {
		parsedLimit, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("invalid limit: %w", err)
		}
		limit = parsedLimit
	}
	user, err := s.db.GetUser(context.Background(), s.Config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("could not find user: %w", err)
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("could not fetch posts: %w", err)
	}
	if len(posts) == 0 {
		fmt.Println("No posts found for this user")
		return nil
	}
	for _, post := range posts {

		fmt.Printf("---%v---\n", post.Title.String)
		fmt.Printf("Link:  %s\n", post.Url)
		fmt.Printf("Content: %v\n\n", post.Description.String)
		if post.PublishedAt.Valid {
			fmt.Printf("Published at: %v\n", post.PublishedAt.Time)
		} else {
			fmt.Println("Published at: unknown")
		}
	}
	return nil
}

func NewCommands() *Commands {
	return &Commands{
		Handlers: make(map[string]func(*State, Command) error),
	}
}

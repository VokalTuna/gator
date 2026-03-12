package main

import (
	"context"
	"fmt"
	"time"

	"github.com/VokalTuna/gator/internal/database"
	"github.com/google/uuid"
)

func handleAddfeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("Usage: %s <name> <url>", cmd.Name)
	}
	name := cmd.Args[0]
	url := cmd.Args[1]
	current_time := time.Now()

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: current_time,
		UpdatedAt: current_time,
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})

	if err != nil {
		return fmt.Errorf("Unable to create feed: %w", err)
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Couldn't create a follow: %w", err)
	}

	fmt.Println("Feed created successfully:")
	printFeed(feed, user)
	fmt.Println()
	fmt.Println("Feed followed successfully: ")
	printFeedFollow(feedFollow.UserName, feedFollow.FeedName)
	fmt.Println("===================================")

	return nil
}

func handlerFeeds(s *state, _ command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Unable to get feeds: %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}

	fmt.Printf("Found %d feeds:\n", len(feeds))
	for _, feed := range feeds {
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("Couldn't get user: %w", err)
		}
		printFeed(feed, user)
		fmt.Println("=========================================")
	}

	return nil
}

func printFeed(feed database.Feed, user database.User) {
	fmt.Printf("* ID:      %s\n", feed.ID)
	fmt.Printf("* Created: %v\n", feed.CreatedAt)
	fmt.Printf("* Updated: %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:    %s\n", feed.Name)
	fmt.Printf("* Url:     %s\n", feed.Url)
	fmt.Printf("* User:  %s\n", user.Name)
}

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/VokalTuna/gator/internal/database"
	"github.com/google/uuid"
)

func handlerAggregator(s *state, cmd command) error {
	res, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("couldn't fetch feed: %w", err)
	}
	fmt.Printf("%+v\n", res)
	return nil
}

func handleAddfeed(s *state, cmd command) error {
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}
	if len(cmd.Args) != 2 {
		return fmt.Errorf("Usage: %s <name> <url>", cmd.Name)
	}
	name := cmd.Args[0]
	url := cmd.Args[1]
	current_time := time.Now()

	registeredFeed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
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

	fmt.Println("Feed created successfully:")
	printFeed(registeredFeed)
	fmt.Println()
	fmt.Println("===================================")

	return nil
}

func printFeed(feed database.Feed) {
	fmt.Printf("* ID:      %s\n", feed.ID)
	fmt.Printf("* Created: %v\n", feed.CreatedAt)
	fmt.Printf("* Updated: %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:    %s\n", feed.Name)
	fmt.Printf("* Url:     %s\n", feed.Url)
	fmt.Printf("* UserID:  %s\n", feed.UserID)
}

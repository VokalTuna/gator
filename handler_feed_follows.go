package main

import (
	"context"
	"fmt"
	"time"

	"github.com/VokalTuna/gator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %s <feed_url>", cmd.Name)
	}
	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("No such feed: %w", err)
	}

	result, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Couldn't create a follow: %w", err)
	}

	printFeedFollow(result.UserName, result.FeedName)
	return nil
}

func handlerListFeedFollows(s *state, cmd command, user database.User) error {
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("No feeds: %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feed follows found for this user.")
		return nil
	}

	fmt.Printf("User %s is following these feeds:\n", user.Name)
	for _, feed_follow := range feeds {
		fmt.Printf("* %s\n", feed_follow.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %s <feed url>", cmd.Name)
	}
	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("No such feed: %w", err)
	}
	err = s.db.DeleteFeedFollows(context.Background(), database.DeleteFeedFollowsParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})

	if err != nil {
		return fmt.Errorf("Unable to delete: %w", err)
	}
	fmt.Println("Feed is no longer followed.")

	return nil
}

func printFeedFollow(username, feedname string) {
	fmt.Printf("* User:          %s\n", username)
	fmt.Printf("* Feed:          %s\n", feedname)
}

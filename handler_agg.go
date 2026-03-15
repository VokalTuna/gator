package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/VokalTuna/gator/internal/database"
	"github.com/google/uuid"
)

func handlerAggregator(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %s <time duration>", cmd.Name)
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("Not a time duration: %w", err)
	}
	fmt.Printf("Collecting feeds every %v\n", timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		fmt.Printf("No feeds to fetch: %v", err)
		return
	}
	_, err = s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		fmt.Printf("Unable to mark fetched feed: %v", err)
		return
	}
	rssfeed, err := fetchFeed(context.Background(), feed.Url)

	if err != nil {
		fmt.Printf("Couldn't collect feed %s: %v", feed.Name, err)
	}

	fmt.Printf("Fetched from %s:\n", rssfeed.Channel.Title)

	for _, item := range rssfeed.Channel.Item {
		addPost(s, item, feed)
	}
}

func addPost(s *state, item RSSItem, feed database.Feed) {
	_, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
		ID:          uuid.New(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Title:       sql.NullString{String: item.Title, Valid: item.Title != ""},
		Url:         item.Link,
		Description: sql.NullString{String: item.Description, Valid: item.Description != ""},
		PublishedAt: convertStringToTime(item.PubDate),
		FeedID:      feed.ID,
	})

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return
		}
		log.Printf("error creating post: %v", err)
	}
}

func convertStringToTime(pubDate string) sql.NullTime {
	if pubDate == "" {
		return sql.NullTime{Valid: false}
	}

	formats := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC3339,
	}

	for _, format := range formats {
		t, err := time.Parse(format, pubDate)
		if err == nil {
			return sql.NullTime{Time: t, Valid: true}
		}
	}
	return sql.NullTime{Valid: false}
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	if len(cmd.Args) > 2 {
		return fmt.Errorf("Usage: %s <int>", cmd.Name)
	}
	limit := 2
	if len(cmd.Args) == 1 {
		if specifiedLimit, err := strconv.Atoi(cmd.Args[0]); err == nil {
			limit = specifiedLimit
		} else {
			return fmt.Errorf("invalid limit: %w", err)
		}
	}
	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})

	if err != nil {
		return fmt.Errorf("Unable to fetch posts: %w", err)
	}
	fmt.Printf("There were found %d posts for %s\n\n", len(posts), user.Name)

	for _, post := range posts {
		fmt.Printf("* %s\n", post.Title.String)
		fmt.Printf("    * %s\n", post.PublishedAt.Time.Format("Mon Jan 1"))
		fmt.Printf("    * %s\n", post.FeedName)
		fmt.Printf("    * %s\n", post.Description.String)
		fmt.Printf("    * %s\n", post.Url)
	}

	return nil
}

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/VokalTuna/gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %s <name>", cmd.Name)
	}
	name := cmd.Args[0]
	user, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("Not an user: %w", err)
	}
	err = s.cfg.SetUser(user.Name)

	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("User switched successfully")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %s <name>", cmd.Name)
	}
	name := cmd.Args[0]
	createdAt := time.Now()
	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
		Name:      name,
	}
	registeredUser, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return fmt.Errorf("User already exists: %w", err)
	}
	s.cfg.SetUser(registeredUser.Name)
	fmt.Printf("User was created. \n %+v\n", registeredUser)
	return nil
}

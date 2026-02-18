package main

import (
	"fmt"

	"github.com/VokalTuna/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("failed to read file: %v\n", err)
	}
	cfg.SetUser("Knut")
	cfg, err = config.Read()
	if err != nil {
		fmt.Printf("failed to read file: %v\n", err)
	}

	fmt.Println(cfg)
}

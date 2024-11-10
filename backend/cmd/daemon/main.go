package main

import (
	"log"
	"os"

	"github.com/charleshuang3/autoget/backend/lib/config"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: daemon <config-path>")
	}

	configPath := os.Args[1]

	_, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %s\n", err)
	}

}

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/charleshuang3/autoget/backend/lib/config"
	"github.com/charleshuang3/autoget/backend/lib/handlers"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: daemon <config-path>")
	}

	configPath := os.Args[1]

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %s\n", err)
	}

	router := handlers.NewRouter(cfg)

	// Handling signals for graceful shutdown.
	stopChan := make(chan struct{})
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChan
		close(stopChan)
	}()

	handlers.StartServer(cfg, router, stopChan)
}

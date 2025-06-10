package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charleshuang3/autoget/backend/internal/config"
	"github.com/charleshuang3/autoget/backend/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func main() {
	configPath := flag.String("c", os.Getenv("CONFIG_PATH"), "path to the configuration file")
	flag.Parse()

	if *configPath == "" {
		log.Fatal().Msg("config path is required")
	}

	cfg, err := config.ReadConfig(*configPath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read config")
	}

	service := handlers.NewService(cfg)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	rg := r.Group("/api/v1")
	service.SetupRouter(rg)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("listen")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exiting")
}

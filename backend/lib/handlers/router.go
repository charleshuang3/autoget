package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/charleshuang3/autoget/backend/lib/config"
	"github.com/charleshuang3/autoget/backend/lib/handlers/prowlarr"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(config *config.Config) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	parr := prowlarr.New(config)

	r.Get("/api/indexers", parr.IndexersHandler)
	return r
}

func StartServer(config *config.Config, router *chi.Mux, stopChan <-chan struct{}) {
	addr := fmt.Sprintf(":%d", config.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		log.Printf("Starting server on %s\n", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", addr, err)
		}
	}()

	<-stopChan
	gracefulShutdown(srv)
}

func gracefulShutdown(srv *http.Server) {
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server gracefully stopped")
}

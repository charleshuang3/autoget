package handlers

import (
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

package api

import (
	"context"
	"log"

	"github.com/OZIOisgood/gamma/internal/db"
	"github.com/OZIOisgood/gamma/internal/storage"
	"github.com/OZIOisgood/gamma/internal/uploads"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	Router *chi.Mux
	Pool   *pgxpool.Pool
}

func NewServer(pool *pgxpool.Pool) *Server {
	s := &Server{
		Router: chi.NewRouter(),
		Pool:   pool,
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(middleware.StripSlashes)

	queries := db.New(s.Pool)

	storageService := storage.New()
	if err := storageService.EnsureBucketExists(context.Background()); err != nil {
		log.Printf("Failed to ensure bucket exists: %v", err)
	}

	uploadsHandler := uploads.NewHandler(storageService, queries)
	uploadsHandler.RegisterRoutes(s.Router)
}

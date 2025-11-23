package api

import (
	"context"
	"log"
	"time"

	"github.com/OZIOisgood/gamma/internal/db"
	"github.com/OZIOisgood/gamma/internal/storage"
	"github.com/OZIOisgood/gamma/internal/uploads"
	"github.com/OZIOisgood/gamma/internal/webhooks"
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
	storageService := s.initStorage()

	uploadsHandler := uploads.NewHandler(storageService, queries)
	uploadsHandler.RegisterRoutes(s.Router)

	webhooksHandler := webhooks.NewHandler(queries)
	webhooksHandler.RegisterRoutes(s.Router)
}

func (s *Server) initStorage() *storage.Storage {
	storageService := storage.New()
	if err := storageService.EnsureBucketExists(context.Background()); err != nil {
		log.Printf("Failed to ensure bucket exists: %v", err)
	}

	// Run notification setup in background
	// MinIO validates the webhook endpoint immediately, so the server must be running first.
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("Failed to configure bucket notification: timeout")
				return
			case <-ticker.C:
				if err := storageService.EnsureBucketNotification(ctx); err == nil {
					log.Println("Bucket notification configured successfully")
					return
				}
			}
		}
	}()

	return storageService
}

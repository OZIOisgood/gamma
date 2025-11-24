package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/OZIOisgood/gamma/internal/auth"
	"github.com/OZIOisgood/gamma/internal/db"
	"github.com/OZIOisgood/gamma/internal/events"
	"github.com/OZIOisgood/gamma/internal/storage"
	"github.com/OZIOisgood/gamma/internal/uploads"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	Router   *chi.Mux
	Pool     *pgxpool.Pool
	EventBus *events.EventBus
	Hub      *Hub
}

func NewServer(pool *pgxpool.Pool, eventBus *events.EventBus) *Server {
	s := &Server{
		Router:   chi.NewRouter(),
		Pool:     pool,
		EventBus: eventBus,
		Hub:      NewHub(),
	}
	go s.Hub.Run()
	s.subscribeToEvents()
	s.routes()
	return s
}

func (s *Server) routes() {
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(middleware.StripSlashes)

	s.Router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	s.Router.Post("/auth/login", auth.Login)
	s.Router.Post("/auth/logout", auth.Logout)
	s.Router.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(s.Hub, w, r)
	})

	queries := db.New(s.Pool)
	storageService := s.initStorage()

	uploadsHandler := uploads.NewHandler(storageService, queries)

	s.Router.Group(func(r chi.Router) {
		r.Use(auth.Middleware)
		uploadsHandler.RegisterRoutes(r)
	})
}

func (s *Server) initStorage() *storage.Storage {
	storageService := storage.New()
	if err := storageService.EnsureBucketExists(context.Background()); err != nil {
		log.Printf("Failed to ensure bucket exists: %v", err)
	}

	if err := storageService.EnsureHLSPublicPolicy(context.Background()); err != nil {
		log.Printf("Failed to ensure HLS public policy: %v", err)
	}

	// Run notification setup in background
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

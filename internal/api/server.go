package api

import (
	"github.com/OZIOisgood/gamma/internal/db"
	"github.com/OZIOisgood/gamma/internal/todos"
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
	todosHandler := todos.NewHandler(queries)
	todosHandler.RegisterRoutes(s.Router)
}

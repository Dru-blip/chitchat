package api

import (
	"chitchat/internal/auth"
	"chitchat/internal/db"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	store *db.Store
	api   *chi.Mux
}

func NewServer(store *db.Store) *Server {
	return &Server{
		store: store,
		api:   chi.NewRouter(),
	}
}

func (s *Server) RegisterRoutes() {
	s.api.Use(middleware.Logger)
	authHandler := auth.NewHandler(s.store)
	s.api.Mount("/auth", authHandler.Register())
}

func (s *Server) Start() {
	http.ListenAndServe(":5000", s.api)
}

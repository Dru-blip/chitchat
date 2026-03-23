package api

import (
	"chitchat/internal/auth"
	"chitchat/internal/db"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type Server struct {
	store *db.Store
	api   *echo.Echo
}

func NewServer(store *db.Store) *Server {
	api := echo.New()
	api.Use(middleware.RequestLogger())
	api.Use(middleware.Recover())
	api.Use(middleware.RequestID())
	return &Server{
		store: store,
		api:   api,
	}
}

func (s *Server) RegisterRoutes() {
	authHandler := auth.NewHandler(s.store.Queries)
	authHandler.Register(s.api)
}

func (s *Server) Start() {
	s.api.Start(":5000")
}

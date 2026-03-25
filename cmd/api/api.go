package api

import (
	"chitchat/internal/auth"
	"chitchat/internal/db"
	"chitchat/internal/utils"

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
	api.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowCredentials: true,
	}))
	api.Validator = utils.NewValidator()
	api.HTTPErrorHandler = utils.GlobalErrorHandler
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
	s.api.Start(":5050")
}

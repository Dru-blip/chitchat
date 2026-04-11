package api

import (
	"chitchat/internal/auth"
	"chitchat/internal/db"
	"chitchat/internal/mailer"
	"chitchat/internal/utils"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type Server struct {
	store  *db.Store
	api    *echo.Echo
	Mailer *mailer.Mailer
}

func NewServer(store *db.Store) (*Server, error) {
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

	Mailer, err := mailer.New()

	if err != nil {
		return nil, err
	}

	return &Server{
		store:  store,
		api:    api,
		Mailer: Mailer,
	}, nil
}

func (s *Server) RegisterRoutes() {
	authService := auth.NewService(s.store.Queries, s.Mailer)
	authHandler := auth.NewHandler(authService, s.api.Logger)
	authHandler.Register(s.api)
}

func (s *Server) Start() {
	s.api.Start(":5050")
}

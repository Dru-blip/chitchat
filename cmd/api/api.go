package api

import (
	"chitchat/internal/auth"
	"chitchat/internal/db"
	"chitchat/internal/mailer"
	"chitchat/internal/utils"
	"encoding/gob"
	"net/http"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type Server struct {
	store          *db.Store
	api            *echo.Echo
	Mailer         *mailer.Mailer
	sessionManager *scs.SessionManager
}

func NewServer(store *db.Store) (*Server, error) {
	api := echo.New()

	//TODO: Move session manager creation into a factory function
	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(store.Db)

	sessionManager.Lifetime = 360 * time.Hour
	sessionManager.Cookie.Name = "chisession"
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.Secure = false

	gob.Register(auth.SessionStore{})

	api.Use(middleware.RequestLogger())
	api.Use(middleware.Recover())
	api.Use(middleware.RequestID())
	api.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowCredentials: true,
	}))
	api.Use(echo.WrapMiddleware(sessionManager.LoadAndSave))
	api.Validator = utils.NewValidator()
	api.HTTPErrorHandler = utils.GlobalErrorHandler

	api.Use(NewSessionMiddleware(sessionManager))

	Mailer, err := mailer.New()

	if err != nil {
		return nil, err
	}

	return &Server{
		store:          store,
		api:            api,
		Mailer:         Mailer,
		sessionManager: sessionManager,
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

func NewSessionMiddleware(store *scs.SessionManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			c.Set("_session", store)
			return next(c)
		}
	}
}

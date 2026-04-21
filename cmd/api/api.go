package api

import (
	"chitchat/internal/auth"
	"chitchat/internal/db"
	"chitchat/internal/keys"
	"chitchat/internal/users"
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
	Mailer         Mailer
	sessionManager *scs.SessionManager
}

func NewServer(store *db.Store, mailer Mailer) (*Server, error) {
	gob.Register(auth.SessionStore{})

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

	api.Use(auth.NewSessionMiddleware(sessionManager))

	return &Server{
		store:          store,
		api:            api,
		Mailer:         mailer,
		sessionManager: sessionManager,
	}, nil
}

func (s *Server) RegisterRoutes() {
	authService := auth.NewService(s.store.Queries, s.Mailer)
	authHandler := auth.NewHandler(authService, s.api.Logger)
	authHandler.Register(s.api)

	usersService := users.NewService(s.store.Queries)
	usersHandler := users.NewHandler(usersService, s.api.Logger)
	usersHandler.Register(s.api)

	keyService := keys.NewService(s.store.Queries)
	keyHandler := keys.NewHandler(keyService, s.api.Logger)
	keyHandler.Register(s.api)
}

func (s *Server) Start() {
	s.api.Start(":5050")
}

func (s *Server) Echo() *echo.Echo {
	return s.api
}

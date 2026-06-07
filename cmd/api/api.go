package api

import (
	"chitchat/internal/auth"
	"chitchat/internal/conversations"
	"chitchat/internal/db"
	"chitchat/internal/keys"
	"chitchat/internal/messages"
	"chitchat/internal/users"
	"chitchat/internal/utils"
	"chitchat/internal/ws"
	"encoding/gob"

	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/redis/go-redis/v9"
)

type App struct {
	store          *db.Store
	api            *echo.Echo
	Mailer         Mailer
	sessionManager *scs.SessionManager
	wsHub          *ws.Hub
	rdb            *redis.Client
}

func NewApp(store *db.Store, mailer Mailer, rdb *redis.Client) (*App, error) {
	gob.Register(auth.SessionStore{})
	sessionManager := auth.NewSessionManager(rdb)
	api := SetupEcho(sessionManager)
	wsHub := ws.NewHub()

	return &App{
		store:          store,
		api:            api,
		Mailer:         mailer,
		sessionManager: sessionManager,
		wsHub:          wsHub,
		rdb:            rdb,
	}, nil
}

func (s *App) RegisterRoutes() {
	authService := auth.NewService(s.store.Queries, s.Mailer)
	authHandler := auth.NewHandler(authService, s.api.Logger, s.rdb)
	authHandler.Register(s.api)

	usersService := users.NewService(s.store.Queries)
	usersHandler := users.NewHandler(usersService, s.api.Logger, s.rdb)
	usersHandler.Register(s.api)

	keyService := keys.NewService(s.store.Queries)
	keyHandler := keys.NewHandler(keyService, s.api.Logger)
	keyHandler.Register(s.api)

	convService := conversations.NewService(s.store.Queries)
	convHandler := conversations.NewHandler(convService, s.wsHub)
	convHandler.Register(s.api)

	msgService := messages.NewService(s.store)
	msgHandler := messages.NewHandler(msgService)
	msgHandler.Register(s.api)

	wsHandler := ws.NewHandler(s.wsHub)
	wsHandler.Register(s.api)
}

func (s *App) Start() {
	s.api.Start(":5050")
}

func (s *App) Echo() *echo.Echo {
	return s.api
}

func SetupEcho(sessionManager *scs.SessionManager) *echo.Echo {
	api := echo.New()

	api.Use(middleware.RequestLogger())
	api.Use(middleware.Recover())
	api.Use(middleware.RequestID())
	api.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowCredentials: true,
	}))

	api.Use(echo.WrapMiddleware(sessionManager.LoadAndSave))
	api.Use(auth.NewSessionMiddleware(sessionManager))

	api.Validator = utils.NewValidator()
	api.HTTPErrorHandler = utils.GlobalErrorHandler

	return api
}

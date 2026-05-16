package api

import (
	"chitchat/internal/auth"
	"chitchat/internal/db"
	"chitchat/internal/keys"
	"chitchat/internal/users"
	"chitchat/internal/utils"
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/goredisstore"
	"github.com/alexedwards/scs/v2"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	store          *db.Store
	api            *echo.Echo
	Mailer         Mailer
	sessionManager *scs.SessionManager
	rdb            *redis.Client
	mqttClient     mqtt.Client
}

func NewServer(store *db.Store, mailer Mailer, rdb *redis.Client) (*Server, error) {
	mqttClient, err := MQTT()
	if err != nil {
		return nil, err
	}
	gob.Register(auth.SessionStore{})

	api := echo.New()

	//TODO: Move session manager creation into a factory function
	sessionManager := scs.New()
	sessionManager.Store = goredisstore.New(rdb)

	sessionManager.Lifetime = 360 * time.Hour
	sessionManager.Cookie.Name = "chisession"
	sessionManager.Cookie.Path = "/"
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.Secure = false

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

	return &Server{
		store:          store,
		api:            api,
		Mailer:         mailer,
		sessionManager: sessionManager,
		rdb:            rdb,
		mqttClient:     mqttClient,
	}, nil
}

func (s *Server) RegisterRoutes() {
	authService := auth.NewService(s.store.Queries, s.Mailer)
	authHandler := auth.NewHandler(authService, s.api.Logger, s.rdb)
	authHandler.Register(s.api)

	usersService := users.NewService(s.store.Queries)
	usersHandler := users.NewHandler(usersService, s.api.Logger, s.rdb)
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

func MQTT() (mqtt.Client, error) {
	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().AddBroker("tcp://broker.emqx.io:1883")
	opts.SetClientID(os.Getenv("EMQX_CLIENTID")).SetPassword(os.Getenv("EMQX_CLIENT_PASSWORD"))

	opts.SetKeepAlive(60 * time.Second)
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return c, nil
}

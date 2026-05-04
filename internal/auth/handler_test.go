package auth_test

import (
	"bytes"
	"chitchat/cmd/api"
	"chitchat/internal/db"
	"chitchat/internal/mailer"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/labstack/echo/v5"
	"github.com/pressly/goose/v3"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	redisContainer "github.com/testcontainers/testcontainers-go/modules/redis"
)

type AuthTestSuite struct {
	suite.Suite
	ctx            context.Context
	ctr            *postgres.PostgresContainer
	redisContainer *redisContainer.RedisContainer
	store          *db.Store
	app            *api.Server
	mailer         *mailer.MockMailer
	rdb            *redis.Client
}

func (s *AuthTestSuite) SetupSuite() {
	ctx := context.Background()

	s.ctx = ctx

	dbName := "users"
	dbUser := "user"
	dbPassword := "password"

	postgresContainer, err := postgres.Run(s.ctx,
		"postgres:18.3-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.BasicWaitStrategies(),
	)

	s.Require().NoError(err)

	rdbContainer, err := redisContainer.Run(s.ctx, "redis:alpine", testcontainers.WithExposedPorts("6379/tcp"))
	s.Require().NoError(err)
	s.redisContainer = rdbContainer

	endpoint, err := rdbContainer.Endpoint(s.ctx, "")
	s.Require().NoError(err)

	rdb := redis.NewClient(&redis.Options{
		Addr:     endpoint,
		Password: "",
		DB:       0,
	})

	s.rdb = rdb

	s.ctr = postgresContainer
	conn, err := postgresContainer.ConnectionString(s.ctx)
	s.Require().NoError(err)

	s.Require().NoError(goose.SetDialect("postgres"))

	sdb, err := sql.Open("pgx", conn)
	s.Require().NoError(err)
	defer sdb.Close()

	s.Require().NoError(goose.Up(sdb, "../db/migrations"))

	store, err := db.Connect(conn)
	s.Require().NoError(err)

	s.store = store

	s.mailer = new(mailer.MockMailer)

	api, err := api.NewServer(s.store, s.mailer, s.rdb)
	s.Require().NoError(err)

	s.app = api

	s.app.RegisterRoutes()
}

func (s *AuthTestSuite) TearDownSuite() {
	s.store.Db.Close()
	s.Require().NoError(s.rdb.Close())
	s.Require().NoError(testcontainers.TerminateContainer(s.ctr))
	s.Require().NoError(testcontainers.TerminateContainer(s.redisContainer))
}

func (s *AuthTestSuite) SetupTest() {
	_, err := s.store.Db.Exec(s.ctx, "TRUNCATE TABLE users,devices,magic_link_sessions RESTART IDENTITY CASCADE;")
	s.Require().NoError(err)

	_, err = s.rdb.FlushAll(s.ctx).Result()
	s.Require().NoError(err)

	s.mailer.Calls = nil
	s.mailer.ExpectedCalls = nil
}

func (s *AuthTestSuite) do(method, path string, body any) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		b, err := json.Marshal(body)
		s.Require().NoError(err)
		reqBody = bytes.NewBuffer(b)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	s.app.Echo().ServeHTTP(rec, req)
	return rec
}

func (s *AuthTestSuite) decodeBody(rec *httptest.ResponseRecorder, v any) {
	s.T().Helper()
	s.Require().NoError(json.NewDecoder(rec.Body).Decode(v))
}

func (s *AuthTestSuite) TestSendMagicLink_InvalidEmail() {
	rec := s.do(http.MethodPost, "/auth/send-magic-link", map[string]any{
		"email":  "pikachu@",
		"pubkey": "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCZKD32iSCQ0d",
	})

	s.Require().Equal(http.StatusUnprocessableEntity, rec.Code)
	s.mailer.AssertNotCalled(s.T(), "SendMagicLink", mock.Anything, mock.Anything)
}

func (s *AuthTestSuite) TestSendMagicLink_Success() {
	s.mailer.On("SendMagicLink", "pikachu@gmail.com", mock.Anything).Return(nil)

	rec := s.do(http.MethodPost, "/auth/send-magic-link", map[string]any{
		"email":  "pikachu@gmail.com",
		"pubkey": "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCZKD32iSCQ0d",
	})

	s.Require().Equal(http.StatusOK, rec.Code)
	s.mailer.AssertExpectations(s.T())
}

func (s *AuthTestSuite) TestSendMagicLink_cooldown() {
	s.mailer.On("SendMagicLink", "pikachu@gmail.com", mock.Anything).Return(nil)

	s.Run("send magic link", func() {
		rec := s.do(http.MethodPost, "/auth/send-magic-link", map[string]any{
			"email":  "pikachu@gmail.com",
			"pubkey": "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCZKD32iSCQ0d",
		})

		s.Require().Equal(http.StatusOK, rec.Code)
	})

	s.Run("cooldown", func() {
		rec := s.do(http.MethodPost, "/auth/send-magic-link", map[string]any{
			"email":  "pikachu@gmail.com",
			"pubkey": "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCZKD32iSCQ0d",
		})

		s.Require().Equal(http.StatusTooManyRequests, rec.Code)
	})

	s.mailer.AssertNumberOfCalls(s.T(), "SendMagicLink", 1)
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}

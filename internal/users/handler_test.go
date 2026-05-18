package users_test

import (
	"bytes"
	"chitchat/cmd/api"
	"chitchat/internal/db"
	"chitchat/internal/mailer"
	"chitchat/internal/utils"
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
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	redisContainer "github.com/testcontainers/testcontainers-go/modules/redis"
)

type UsersTestSuite struct {
	suite.Suite
	ctx            context.Context
	ctr            *postgres.PostgresContainer
	redisContainer *redisContainer.RedisContainer
	store          *db.Store
	app            *api.App
	mailer         *mailer.MockMailer
	rdb            *redis.Client
}

func (s *UsersTestSuite) SetupSuite() {
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

	apiServer, err := api.NewApp(s.store, s.mailer, s.rdb)
	s.Require().NoError(err)
	s.app = apiServer
	s.app.RegisterRoutes()
}

func (s *UsersTestSuite) TearDownSuite() {
	s.store.Db.Close()
	s.Require().NoError(s.rdb.Close())
	s.Require().NoError(testcontainers.TerminateContainer(s.ctr))
	s.Require().NoError(testcontainers.TerminateContainer(s.redisContainer))
}

func (s *UsersTestSuite) SetupTest() {
	_, err := s.store.Db.Exec(s.ctx, "TRUNCATE TABLE users,devices,magic_link_sessions,device_signed_prekeys,device_prekeys RESTART IDENTITY CASCADE;")
	s.Require().NoError(err)

	_, err = s.rdb.FlushAll(s.ctx).Result()
	s.Require().NoError(err)

	s.mailer.Calls = nil
	s.mailer.ExpectedCalls = nil
}

func (s *UsersTestSuite) do(method, path string, body any) *httptest.ResponseRecorder {
	return s.doWithCookies(method, path, body)
}

func (s *UsersTestSuite) doWithCookies(method, path string, body any, cookies ...*http.Cookie) *httptest.ResponseRecorder {
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
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	rec := httptest.NewRecorder()
	s.app.Echo().ServeHTTP(rec, req)
	return rec
}

func (s *UsersTestSuite) decodeBody(rec *httptest.ResponseRecorder, v any) {
	s.T().Helper()
	s.Require().NoError(json.NewDecoder(rec.Body).Decode(v))
}

func (s *UsersTestSuite) seedMagicLinkSession(email, pubkey, rawToken string) {
	s.T().Helper()
	_, err := s.store.Db.Exec(s.ctx, `
		INSERT INTO magic_link_sessions
			(email, pubkey, token, ip_address, expires_at, attempts, status)
		VALUES
			($1, $2, $3, '127.0.0.1', NOW() + INTERVAL '15 minutes', 1, 'pending')
	`, email, pubkey, utils.SHA256(rawToken))
	s.Require().NoError(err)
}

func (s *UsersTestSuite) loginAndReturnSession(email, pubkey, rawToken string) (*http.Cookie, string) {
	s.T().Helper()

	s.seedMagicLinkSession(email, pubkey, rawToken)

	rec := s.do(http.MethodPost, "/auth/verify-magic-link", map[string]any{
		"token": rawToken,
	})
	s.Require().Equal(http.StatusOK, rec.Code)

	var res map[string]any
	s.decodeBody(rec, &res)

	for _, cookie := range rec.Result().Cookies() {
		if cookie.Name == "chisession" {
			return cookie, res["userId"].(string)
		}
	}

	s.FailNow("session cookie not found")
	return nil, ""
}

// ─── Onboard ──────────────────────────────────────────────────────────────────

func (s *UsersTestSuite) TestOnboard_RequiresAuth() {
	rec := s.do(http.MethodPatch, "/users/onboard", map[string]any{
		"name":     "Ash Ketchum",
		"pubkey":   "device-pubkey",
		"password": "pikachu123",
		"image":    "https://example.com/ash.png",
	})
	s.Require().Equal(http.StatusUnauthorized, rec.Code)
}

func (s *UsersTestSuite) TestOnboard_InvalidPayload_MissingName() {
	sessionCookie, _ := s.loginAndReturnSession("ash@pallet.com", "device-pubkey-ash", "token-ash")

	rec := s.doWithCookies(http.MethodPatch, "/users/onboard", map[string]any{
		"pubkey":   "device-pubkey-ash",
		"password": "pikachu123",
	}, sessionCookie)

	s.Require().Equal(http.StatusUnprocessableEntity, rec.Code)
}

func (s *UsersTestSuite) TestOnboard_InvalidPayload_MissingPassword() {
	sessionCookie, _ := s.loginAndReturnSession("misty@cerulean.com", "device-pubkey-misty", "token-misty")

	rec := s.doWithCookies(http.MethodPatch, "/users/onboard", map[string]any{
		"name":   "Misty",
		"pubkey": "device-pubkey-misty",
	}, sessionCookie)

	s.Require().Equal(http.StatusUnprocessableEntity, rec.Code)
}

func (s *UsersTestSuite) TestOnboard_InvalidPayload_MissingPubkey() {
	sessionCookie, _ := s.loginAndReturnSession("brock@pewter.com", "device-pubkey-brock", "token-brock")

	rec := s.doWithCookies(http.MethodPatch, "/users/onboard", map[string]any{
		"name":     "Brock",
		"password": "rock123",
	}, sessionCookie)

	s.Require().Equal(http.StatusUnprocessableEntity, rec.Code)
}

func (s *UsersTestSuite) TestOnboard_InvalidPayload_EmptyBody() {
	sessionCookie, _ := s.loginAndReturnSession("gary@pallet.com", "device-pubkey-gary", "token-gary")

	rec := s.doWithCookies(http.MethodPatch, "/users/onboard", map[string]any{}, sessionCookie)

	s.Require().Equal(http.StatusUnprocessableEntity, rec.Code)
}

func (s *UsersTestSuite) TestOnboard_Success() {
	email := "ash@pallet.com"
	pubkey := "device-pubkey-ash"
	sessionCookie, userID := s.loginAndReturnSession(email, pubkey, "token-ash")

	rec := s.doWithCookies(http.MethodPatch, "/users/onboard", map[string]any{
		"name":     "Ash Ketchum",
		"pubkey":   pubkey,
		"password": "pikachu123",
		"image":    "https://example.com/ash.png",
	}, sessionCookie)

	s.Require().Equal(http.StatusOK, rec.Code)

	var name string
	var image string
	var onboarding bool
	err := s.store.Db.QueryRow(s.ctx,
		"SELECT name, image, onboarding FROM users WHERE id = $1", userID,
	).Scan(&name, &image, &onboarding)
	s.Require().NoError(err)
	s.Require().Equal("Ash Ketchum", name)
	s.Require().Equal("https://example.com/ash.png", image)
	s.Require().False(onboarding)

	var passwordHash string
	err = s.store.Db.QueryRow(s.ctx,
		"SELECT password FROM users WHERE id = $1", userID,
	).Scan(&passwordHash)
	s.Require().NoError(err)
	s.Require().NotEmpty(passwordHash)
	s.Require().NotEqual("pikachu123", passwordHash)
}

func TestUsersSuite(t *testing.T) {
	suite.Run(t, new(UsersTestSuite))
}

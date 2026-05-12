package keys_test

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

type KeysTestSuite struct {
	suite.Suite
	ctx            context.Context
	ctr            *postgres.PostgresContainer
	redisContainer *redisContainer.RedisContainer
	store          *db.Store
	app            *api.Server
	mailer         *mailer.MockMailer
	rdb            *redis.Client
}

func (s *KeysTestSuite) SetupSuite() {
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

	apiServer, err := api.NewServer(s.store, s.mailer, s.rdb)
	s.Require().NoError(err)
	s.app = apiServer
	s.app.RegisterRoutes()
}

func (s *KeysTestSuite) TearDownSuite() {
	s.store.Db.Close()
	s.Require().NoError(s.rdb.Close())
	s.Require().NoError(testcontainers.TerminateContainer(s.ctr))
	s.Require().NoError(testcontainers.TerminateContainer(s.redisContainer))
}

func (s *KeysTestSuite) SetupTest() {
	_, err := s.store.Db.Exec(s.ctx, "TRUNCATE TABLE users,devices,magic_link_sessions,device_signed_prekeys,device_prekeys RESTART IDENTITY CASCADE;")
	s.Require().NoError(err)

	_, err = s.rdb.FlushAll(s.ctx).Result()
	s.Require().NoError(err)

	s.mailer.Calls = nil
	s.mailer.ExpectedCalls = nil
}

func (s *KeysTestSuite) do(method, path string, body any) *httptest.ResponseRecorder {
	return s.doWithCookies(method, path, body)
}

func (s *KeysTestSuite) doWithCookies(method, path string, body any, cookies ...*http.Cookie) *httptest.ResponseRecorder {
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

func (s *KeysTestSuite) decodeBody(rec *httptest.ResponseRecorder, v any) {
	s.T().Helper()
	s.Require().NoError(json.NewDecoder(rec.Body).Decode(v))
}

func (s *KeysTestSuite) seedMagicLinkSession(email, pubkey, rawToken string) {
	s.T().Helper()
	_, err := s.store.Db.Exec(s.ctx, `
        INSERT INTO magic_link_sessions
            (email, pubkey, token, ip_address, expires_at, attempts, status)
        VALUES
            ($1, $2, $3, '127.0.0.1', NOW() + INTERVAL '15 minutes', 1, 'pending')
    `, email, pubkey, utils.SHA256(rawToken))
	s.Require().NoError(err)
}

func (s *KeysTestSuite) loginAndReturnSession(email, pubkey, rawToken string) (*http.Cookie, string) {
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

// ─── Upload Prekeys ───────────────────────────────────────────────────────────

func (s *KeysTestSuite) TestUploadPrekeys_RequiresAuth() {
	rec := s.do(http.MethodPost, "/keys/", map[string]any{
		"prekeyIds": []int{1},
		"prekeys":   []string{"prekey-1"},
		"signedPreKey": map[string]any{
			"id":        1,
			"key":       "signed-key",
			"signature": "signature",
		},
	})
	s.Require().Equal(http.StatusUnauthorized, rec.Code)
}

func (s *KeysTestSuite) TestUploadPrekeys_StoresKeysAndAllowsDuplicateRetry() {
	sessionCookie, _ := s.loginAndReturnSession("brock@pewter.com", "device-pubkey", "token-brock")

	payload := map[string]any{
		"prekeyIds": []int{1, 2},
		"prekeys":   []string{"prekey-1", "prekey-2"},
		"signedPreKey": map[string]any{
			"id":        1,
			"key":       "signed-key",
			"signature": "signature",
		},
	}

	rec := s.doWithCookies(http.MethodPost, "/keys/", payload, sessionCookie)
	s.Require().Equal(http.StatusOK, rec.Code)

	// duplicate upload should also succeed (upsert behaviour)
	rec = s.doWithCookies(http.MethodPost, "/keys/", payload, sessionCookie)
	s.Require().Equal(http.StatusOK, rec.Code)

	var signedCount int
	err := s.store.Db.QueryRow(s.ctx, "SELECT COUNT(*) FROM device_signed_prekeys WHERE public_key=$1", "signed-key").Scan(&signedCount)
	s.Require().NoError(err)
	s.Require().Equal(1, signedCount)

	var prekeyCount int
	err = s.store.Db.QueryRow(s.ctx, "SELECT COUNT(*) FROM device_prekeys WHERE public_key IN ($1, $2)", "prekey-1", "prekey-2").Scan(&prekeyCount)
	s.Require().NoError(err)
	s.Require().Equal(2, prekeyCount)
}

func (s *KeysTestSuite) TestUploadPrekeys_RejectsMismatchedPrekeyIdsAndKeys() {
	sessionCookie, _ := s.loginAndReturnSession("erika@celadon.com", "device-pubkey-erika", "token-erika")

	rec := s.doWithCookies(http.MethodPost, "/keys/", map[string]any{
		"prekeyIds": []int{1, 2},
		"prekeys":   []string{"prekey-1"},
		"signedPreKey": map[string]any{
			"id":        1,
			"key":       "signed-key",
			"signature": "signature",
		},
	}, sessionCookie)

	s.Require().Equal(http.StatusBadRequest, rec.Code)
}

func (s *KeysTestSuite) TestUploadPrekeys_InvalidPayload_EmptyPrekeys() {
	sessionCookie, _ := s.loginAndReturnSession("misty@cerulean.com", "device-pubkey-misty", "token-misty")

	rec := s.doWithCookies(http.MethodPost, "/keys/", map[string]any{
		"prekeyIds": []int{},
		"prekeys":   []string{},
		"signedPreKey": map[string]any{
			"id":        1,
			"key":       "signed-key",
			"signature": "signature",
		},
	}, sessionCookie)

	s.Require().Equal(http.StatusUnprocessableEntity, rec.Code)
}

func (s *KeysTestSuite) TestUploadPrekeys_InvalidPayload_MissingSignedPreKey() {
	sessionCookie, _ := s.loginAndReturnSession("ash@pallet.com", "device-pubkey-ash", "token-ash")

	rec := s.doWithCookies(http.MethodPost, "/keys/", map[string]any{
		"prekeyIds": []int{1},
		"prekeys":   []string{"prekey-1"},
	}, sessionCookie)

	s.Require().Equal(http.StatusUnprocessableEntity, rec.Code)
}

// ─── Get Key Bundle ───────────────────────────────────────────────────────────

func (s *KeysTestSuite) TestGetKeyBundle_RequiresAuth() {
	rec := s.do(http.MethodPost, "/keys/some-user-id", nil)
	s.Require().Equal(http.StatusUnauthorized, rec.Code)
}

func (s *KeysTestSuite) TestGetKeyBundle_ConsumesOnePrekey() {
	sessionCookie, userID := s.loginAndReturnSession("sabrina@saffron.com", "device-pubkey-sabrina", "token-sabrina")

	upload := map[string]any{
		"prekeyIds": []int{10, 11},
		"prekeys":   []string{"prekey-10", "prekey-11"},
		"signedPreKey": map[string]any{
			"id":        7,
			"key":       "signed-key-7",
			"signature": "signature-7",
		},
	}

	rec := s.doWithCookies(http.MethodPost, "/keys/", upload, sessionCookie)
	s.Require().Equal(http.StatusOK, rec.Code)

	rec = s.doWithCookies(http.MethodPost, "/keys/"+userID, nil, sessionCookie)
	s.Require().Equal(http.StatusOK, rec.Code)

	var remaining int
	err := s.store.Db.QueryRow(s.ctx, "SELECT COUNT(*) FROM device_prekeys").Scan(&remaining)
	s.Require().NoError(err)
	s.Require().Equal(1, remaining)

	var body struct {
		Bundle struct {
			SignedKeyID  int32  `json:"SignedKeyID"`
			SignedPubkey string `json:"SignedPubkey"`
			PrekeyID     int32  `json:"PrekeyID"`
			Prekey       string `json:"Prekey"`
		} `json:"bundle"`
	}
	s.decodeBody(rec, &body)
	s.Require().Equal(int32(7), body.Bundle.SignedKeyID)
	s.Require().Equal("signed-key-7", body.Bundle.SignedPubkey)
	s.Require().Contains([]int32{10, 11}, body.Bundle.PrekeyID)
	s.Require().Contains([]string{"prekey-10", "prekey-11"}, body.Bundle.Prekey)
}

func TestKeysSuite(t *testing.T) {
	suite.Run(t, new(KeysTestSuite))
}

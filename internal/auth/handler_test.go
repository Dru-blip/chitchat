package auth_test

import (
	"bytes"
	"chitchat/cmd/api"
	"chitchat/internal/auth"
	"chitchat/internal/db"
	"chitchat/internal/mailer"
	"chitchat/internal/utils"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
	_, err := s.store.Db.Exec(s.ctx, "TRUNCATE TABLE users,devices,magic_link_sessions,device_signed_prekeys,device_prekeys RESTART IDENTITY CASCADE;")
	s.Require().NoError(err)

	_, err = s.rdb.FlushAll(s.ctx).Result()
	s.Require().NoError(err)

	s.mailer.Calls = nil
	s.mailer.ExpectedCalls = nil
}

func (s *AuthTestSuite) do(method, path string, body any) *httptest.ResponseRecorder {
	return s.doWithCookies(method, path, body)
}

func (s *AuthTestSuite) doWithCookies(method, path string, body any, cookies ...*http.Cookie) *httptest.ResponseRecorder {
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

func (s *AuthTestSuite) decodeBody(rec *httptest.ResponseRecorder, v any) {
	s.T().Helper()
	s.Require().NoError(json.NewDecoder(rec.Body).Decode(v))
}

func (s *AuthTestSuite) seedMagicLinkSession(email, pubkey, rawToken string) {
	s.T().Helper()
	_, err := s.store.Db.Exec(s.ctx, `
		INSERT INTO magic_link_sessions
			(email, pubkey, token, ip_address, expires_at, attempts, status)
		VALUES
			($1, $2, $3, '127.0.0.1', NOW() + INTERVAL '15 minutes', 1, 'pending')
	`, email, pubkey, utils.SHA256(rawToken))
	s.Require().NoError(err)
}

func (s *AuthTestSuite) seedMaxAttemptsSession(email, pubkey string) {
	s.T().Helper()
	_, err := s.store.Db.Exec(s.ctx, `
		INSERT INTO magic_link_sessions
			(email, pubkey, token, ip_address, expires_at, attempts, status)
		VALUES
			($1, $2, 'exhaused-token-hash', '127.0.0.1', NOW() + INTERVAL '15 minutes', $3, 'pending')
	`, email, pubkey, auth.MaxAttempts)
	s.Require().NoError(err)
}

func (s *AuthTestSuite) seedExpiredMagicLinkSession(email, pubkey, rawToken string) {
	s.T().Helper()
	_, err := s.store.Db.Exec(s.ctx, `
		INSERT INTO magic_link_sessions
			(email, pubkey, token, ip_address, expires_at, attempts, status)
		VALUES
			($1, $2, $3, '127.0.0.1', NOW() - INTERVAL '1 hour', 1, 'pending')
	`, email, pubkey, utils.SHA256(rawToken))
	s.Require().NoError(err)
}

func (s *AuthTestSuite) seedUsedMagicLinkSession(email, pubkey, rawToken string) {
	s.T().Helper()
	_, err := s.store.Db.Exec(s.ctx, `
		INSERT INTO magic_link_sessions
			(email, pubkey, token, ip_address, expires_at, attempts, status)
		VALUES
			($1, $2, $3, '127.0.0.1', NOW() + INTERVAL '15 minutes', 1, 'used')
	`, email, pubkey, utils.SHA256(rawToken))
	s.Require().NoError(err)
}

func (s *AuthTestSuite) loginAndReturnSession(email, pubkey, rawToken string) (*http.Cookie, string) {
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

// Validations
func (s *AuthTestSuite) TestSendMagicLink_InvalidEmail() {
	rec := s.do(http.MethodPost, "/auth/send-magic-link", map[string]any{
		"email":  "pikachu@",
		"pubkey": "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCZKD32iSCQ0d",
	})

	s.Require().Equal(http.StatusUnprocessableEntity, rec.Code)
	s.mailer.AssertNotCalled(s.T(), "SendMagicLink", mock.Anything, mock.Anything)
}

func (s *AuthTestSuite) TestSendMagicLink_MissingEmail() {
	rec := s.do(http.MethodPost, "/auth/send-magic-link", map[string]any{
		"pubkey": "somevalidpubkey",
	})

	s.Require().Equal(http.StatusUnprocessableEntity, rec.Code)
	s.mailer.AssertNotCalled(s.T(), "SendMagicLink", mock.Anything, mock.Anything)
}

func (s *AuthTestSuite) TestSendMagicLink_MissingPubkey() {
	rec := s.do(http.MethodPost, "/auth/send-magic-link", map[string]any{
		"email": "pikachu@gmail.com",
	})

	s.Require().Equal(http.StatusUnprocessableEntity, rec.Code)
	s.mailer.AssertNotCalled(s.T(), "SendMagicLink", mock.Anything, mock.Anything)
}

func (s *AuthTestSuite) TestSendMagicLink_EmptyBody() {
	rec := s.do(http.MethodPost, "/auth/send-magic-link", map[string]any{})
	s.Require().Equal(http.StatusUnprocessableEntity, rec.Code)
}

func (s *AuthTestSuite) TestSendMagicLink_SuccessAndResendCoolDown() {
	s.mailer.On("SendMagicLink", "pikachu@gmail.com", mock.Anything).Return(nil)

	s.Run("send magic link", func() {
		rec := s.do(http.MethodPost, "/auth/send-magic-link", map[string]any{
			"email":  "pikachu@gmail.com",
			"pubkey": "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCZKD32iSCQ0d",
		})

		s.Require().Equal(http.StatusOK, rec.Code)
		s.mailer.AssertNumberOfCalls(s.T(), "SendMagicLink", 1)
		s.mailer.AssertExpectations(s.T())

		var res auth.SendMagicLinkResponse
		s.decodeBody(rec, &res)

		s.Require().Equal(res.Email, "pikachu@gmail.com")
		s.Require().Greater(res.RetryAfter, time.Now())
	})

	s.mailer.Calls = nil
	s.mailer.ExpectedCalls = nil

	s.Run("resend cool down", func() {
		rec := s.do(http.MethodPost, "/auth/send-magic-link", map[string]any{
			"email":  "pikachu@gmail.com",
			"pubkey": "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCZKD32iSCQ0d",
		})

		s.Require().Equal(http.StatusTooManyRequests, rec.Code)
		var res auth.SendMagicLinkResponse
		s.decodeBody(rec, &res)

		s.Require().Greater(res.RetryAfter, time.Now())
		s.mailer.AssertNotCalled(s.T(), "SendMagicLink", mock.Anything, mock.Anything)
	})

}

func (s *AuthTestSuite) TestSendMagicLink_DifferentEmails_Independent() {
	s.mailer.On("SendMagicLink", mock.Anything, mock.Anything).Return(nil)

	rec1 := s.do(http.MethodPost, "/auth/send-magic-link", map[string]any{
		"email":  "ash@pallet.com",
		"pubkey": "pubkey1",
	})
	rec2 := s.do(http.MethodPost, "/auth/send-magic-link", map[string]any{
		"email":  "misty@cerulean.com",
		"pubkey": "pubkey2",
	})

	s.Require().Equal(http.StatusOK, rec1.Code)
	s.Require().Equal(http.StatusOK, rec2.Code)
	s.mailer.AssertNumberOfCalls(s.T(), "SendMagicLink", 2)
}

func (s *AuthTestSuite) TestSendMagicLink_MaxAttempts_Returns429WithEmptyRetryAfter() {
	email := "pikachu@gmail.com"
	pubkey := "somevalidpubkey"

	s.seedMaxAttemptsSession(email, pubkey)

	rec := s.do(http.MethodPost, "/auth/send-magic-link", map[string]any{
		"email":  email,
		"pubkey": pubkey,
	})

	s.Require().Equal(http.StatusTooManyRequests, rec.Code)

	var res auth.SendMagicLinkResponse
	s.decodeBody(rec, &res)
	s.Require().True(res.RetryAfter.IsZero())

	s.mailer.AssertNotCalled(s.T(), "SendMagicLink", mock.Anything, mock.Anything)
}

func (s *AuthTestSuite) TestVerifyMagicLink_MissingToken() {
	rec := s.do(http.MethodPost, "/auth/verify-magic-link", map[string]any{})
	s.Require().Equal(http.StatusUnprocessableEntity, rec.Code)
}

func (s *AuthTestSuite) TestVerifyMagicLink_NewUser_CreatesUserAndDevice() {
	email := "ash@pallet.com"
	pubkey := "example"
	rawToken := "fresh-token"

	s.seedMagicLinkSession(email, pubkey, rawToken)

	rec := s.do(http.MethodPost, "/auth/verify-magic-link", map[string]any{
		"token": rawToken,
	})
	s.Require().Equal(http.StatusOK, rec.Code)

	var res map[string]any
	s.decodeBody(rec, &res)

	s.Require().True(res["onboard"].(bool))
	s.Require().True(res["redirect"].(bool))
	s.Require().NotEmpty(res["userId"])

	var userCount int
	err := s.store.Db.QueryRow(s.ctx, "SELECT COUNT(*) FROM users WHERE email=$1", email).Scan(&userCount)
	s.Require().NoError(err)
	s.Require().Equal(1, userCount)

	var deviceCount int
	err = s.store.Db.QueryRow(s.ctx, "SELECT COUNT(*) FROM devices WHERE pubkey=$1", pubkey).Scan(&deviceCount)
	s.Require().NoError(err)
	s.Require().Equal(1, deviceCount)
}

func (s *AuthTestSuite) TestVerifyMagicLink_ExistingUser_SkipsOnboarding() {
	email := "misty@cerulean.com"
	pubkey := "pubkey-misty"
	rawToken := "token-misty"

	_, err := s.store.Db.Exec(s.ctx,
		"INSERT INTO users (email, ipkey, onboarding) VALUES ($1, $2, $3)", email, pubkey, false)
	s.Require().NoError(err)

	s.seedMagicLinkSession(email, pubkey, rawToken)

	rec := s.do(http.MethodPost, "/auth/verify-magic-link", map[string]any{
		"token": rawToken,
	})

	// TODO: test should check for device creation
	s.Require().Equal(http.StatusOK, rec.Code)

	var res map[string]any
	s.decodeBody(rec, &res)
	s.Require().False(res["onboard"].(bool))
}

func (s *AuthTestSuite) TestVerifyMagicLink_MarksTokenAsUsed() {
	email := "thanos@gmail.com"
	pubkey := "pubkey"
	rawToken := "token"

	s.seedMagicLinkSession(email, pubkey, rawToken)

	rec := s.do(http.MethodPost, "/auth/verify-magic-link", map[string]any{
		"token": rawToken,
	})
	s.Require().Equal(http.StatusOK, rec.Code)

	var status string
	err := s.store.Db.QueryRow(s.ctx,
		"SELECT status FROM magic_link_sessions WHERE token=$1",
		utils.SHA256(rawToken),
	).Scan(&status)
	s.Require().NoError(err)
	s.Require().Equal("used", status)
}

func (s *AuthTestSuite) TestVerifyMagicLink_InvalidToken_Returns401() {
	rec := s.do(http.MethodPost, "/auth/verify-magic-link", map[string]any{
		"token": "this-token-was-never-issued",
	})
	s.Require().Equal(http.StatusUnauthorized, rec.Code)
}

func (s *AuthTestSuite) TestVerifyMagicLink_ExpiredToken_Returns401() {
	email := "koga@fuchsia.com"
	pubkey := "pubkey-koga"
	rawToken := "expired-token-koga"

	s.seedExpiredMagicLinkSession(email, pubkey, rawToken)

	rec := s.do(http.MethodPost, "/auth/verify-magic-link", map[string]any{
		"token": rawToken,
	})
	s.Require().Equal(http.StatusUnauthorized, rec.Code)
}

func (s *AuthTestSuite) TestVerifyMagicLink_AlreadyUsedToken_Returns401() {
	email := "blaine@cinnabar.com"
	pubkey := "pubkey-blaine"
	rawToken := "used-token-blaine"

	s.seedUsedMagicLinkSession(email, pubkey, rawToken)

	rec := s.do(http.MethodPost, "/auth/verify-magic-link", map[string]any{
		"token": rawToken,
	})
	s.Require().Equal(http.StatusUnauthorized, rec.Code)
}

func (s *AuthTestSuite) TestVerifyMagicLink_ReplayAttack_SecondUseRejected() {
	email := "giovanni@viridian.com"
	pubkey := "pubkey-giovanni"
	rawToken := "token-giovanni"

	s.seedMagicLinkSession(email, pubkey, rawToken)

	rec := s.do(http.MethodPost, "/auth/verify-magic-link", map[string]any{"token": rawToken})
	s.Require().Equal(http.StatusOK, rec.Code)

	rec = s.do(http.MethodPost, "/auth/verify-magic-link", map[string]any{"token": rawToken})
	s.Require().Equal(http.StatusUnauthorized, rec.Code)
}

func (s *AuthTestSuite) TestMeRequiresAuth() {
	rec := s.do(http.MethodGet, "/auth/me", nil)
	s.Require().Equal(http.StatusUnauthorized, rec.Code)
}

func (s *AuthTestSuite) TestMeReturnsSessionWhenAuthenticated() {
	sessionCookie, userID := s.loginAndReturnSession("janine@fuchsia.com", "device-pubkey-janine", "token-janine")

	rec := s.doWithCookies(http.MethodGet, "/auth/me", nil, sessionCookie)
	s.Require().Equal(http.StatusOK, rec.Code)

	var body auth.SessionStore
	s.decodeBody(rec, &body)
	s.Require().Equal(userID, body.UserId)
	s.Require().Equal("janine@fuchsia.com", body.Email)
}

func (s *AuthTestSuite) TestUploadPrekeys_RequiresAuth() {
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

func (s *AuthTestSuite) TestUploadPrekeys_StoresKeysAndAllowsDuplicateRetry() {
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

func (s *AuthTestSuite) TestUploadPrekeys_RejectsMismatchedPrekeyIdsAndKeys() {
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

func (s *AuthTestSuite) TestGetKeyBundle_ConsumesOnePrekey() {
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

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}

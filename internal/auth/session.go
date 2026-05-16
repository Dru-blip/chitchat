package auth

import (
	"net/http"
	"time"

	"github.com/alexedwards/scs/goredisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/redis/go-redis/v9"
)

func NewSessionManager(rdb *redis.Client) *scs.SessionManager {
	sessionManager := scs.New()
	sessionManager.Store = goredisstore.New(rdb)

	sessionManager.Lifetime = 360 * time.Hour
	sessionManager.Cookie.Name = "chisession"
	sessionManager.Cookie.Path = "/"
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.Secure = false
	return sessionManager
}

package auth

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v5"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		store := c.Get("_session").(*scs.SessionManager)

		user, ok := store.Get(c.Request().Context(), "user").(SessionStore)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "Not Authenticated")
		}
		c.Set("user", user)
		return next(c)
	}
}

func NewSessionMiddleware(store *scs.SessionManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			c.Set("_session", store)
			return next(c)
		}
	}
}

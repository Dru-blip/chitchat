package auth

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

var (
	// General errors
	ErrInvalidRequest = echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	ErrUnauthorized   = echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	ErrInternal       = echo.NewHTTPError(http.StatusInternalServerError, "internal error")

	// Magic Link errors
	ErrInvalidMagicLink = echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired magic link")
	ErrMagicLinkUsed    = echo.NewHTTPError(http.StatusUnauthorized, "magic link already used")
	ErrMagicLinkRevoked = echo.NewHTTPError(http.StatusUnauthorized, "magic link revoked")
	ErrMagicLinkExpired = echo.NewHTTPError(http.StatusUnauthorized, "magic link expired")

	// Session errors
	ErrInvalidSession = echo.NewHTTPError(http.StatusUnauthorized, "session expired or invalid")
	ErrSessionRevoked = echo.NewHTTPError(http.StatusUnauthorized, "session revoked")
	ErrNoSession      = echo.NewHTTPError(http.StatusUnauthorized, "no active session")

	ErrUserNotFound = echo.NewHTTPError(http.StatusNotFound, "user not found")

	ErrTooManyAttempts = echo.NewHTTPError(http.StatusTooManyRequests, "too many attempts")
)

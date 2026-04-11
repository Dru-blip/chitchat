package auth

import "errors"

var (
	// General errors
	ErrInvalidRequest = errors.New("invalid request")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrInternal       = errors.New("internal error")

	// Magic Link errors
	ErrInvalidMagicLink = errors.New("invalid or expired magic link")
	ErrMagicLinkUsed    = errors.New("magic link already used")
	ErrMagicLinkRevoked = errors.New("magic link revoked")
	ErrMagicLinkExpired = errors.New("magic link expired")

	// Session errors
	ErrInvalidSession = errors.New("session expired or invalid")
	ErrSessionRevoked = errors.New("session revoked")
	ErrNoSession      = errors.New("no active session")
)

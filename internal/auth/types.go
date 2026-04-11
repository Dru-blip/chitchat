package auth

import (
	"chitchat/internal/db/sqlc"
	"context"
	"crypto/rand"
	"encoding/base64"

	"github.com/google/uuid"
)

type SendMagicLinkPayload struct {
	Email  string `json:"email" validate:"required,email"`
	Pubkey string `json:"pubkey" validate:"required"`
}

type VerifyMagicLinkPayload struct {
	Token string `json:"token" validate:"required"`
}

type SendMagicLinkResponse struct {
	Message string `json:"message"`
	Email   string `json:"email"`
}

type SessionResponse struct {
	UserID    uuid.UUID `json:"userId"`
	Email     string    `json:"email"`
	ExpiresAt int64     `json:"expiresAt"`
}

type Repository interface {
	CreateMagicLinkSession(ctx context.Context, arg sqlc.CreateMagicLinkSessionParams) (sqlc.MagicLinkSession, error)
	GetMagicLinkSessionByToken(ctx context.Context, token string) (sqlc.MagicLinkSession, error)
	GetPendingMagicLinkSession(ctx context.Context, email string) (sqlc.MagicLinkSession, error)
	MarkMagicLinkAsUsed(ctx context.Context, token string) (sqlc.MagicLinkSession, error)
	RevokeMagicLink(ctx context.Context, token string) error
}

type Mailer interface {
	SendMagicLink(recipient string, link string) error
}

type ContextKey string

const (
	ContextKeyUserID    ContextKey = "userID"
	ContextKeySessionID ContextKey = "sessionID"
	ContextKeyEmail     ContextKey = "email"
)

func GenerateMagicLinkToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

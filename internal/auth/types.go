package auth

import (
	"chitchat/internal/db/sqlc"
	"context"
	"net/netip"
	"time"

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
	Message    string    `json:"message"`
	Email      string    `json:"email"`
	RetryAfter time.Time `json:"retryAfter"`
}

type SessionResponse struct {
	UserID    uuid.UUID `json:"userId"`
	Email     string    `json:"email"`
	ExpiresAt int64     `json:"expiresAt"`
}

type Service interface {
	SendMagicLink(ctx context.Context, email, pubkey string, ipAddress netip.Addr, userAgent string) (*SendMagicLinkResponse, error)
	VerifyMagicLink(ctx context.Context, token string, ipAddress netip.Addr, userAgent string) (*sqlc.MagicLinkSession, error)
	GetOrCreateUser(ctx context.Context, email, pubkey string) (*sqlc.User, bool, error)
	GetOrCreateDevice(ctx context.Context, user_id uuid.UUID, pubkey, os, user_agent string) (*sqlc.Device, error)
}

type SessionInfo struct {
	SessionID string
	UserID    uuid.UUID
	Email     string
	ExpiresAt time.Time
}

type Repository interface {
	CreateMagicLinkSession(ctx context.Context, arg sqlc.CreateMagicLinkSessionParams) (sqlc.MagicLinkSession, error)
	GetMagicLinkSessionByToken(ctx context.Context, token string) (sqlc.MagicLinkSession, error)
	GetPendingMagicLinkSession(ctx context.Context, email string) (sqlc.MagicLinkSession, error)
	MarkMagicLinkAsUsed(ctx context.Context, token string) (sqlc.MagicLinkSession, error)
	UpdateMagicLinkSession(ctx context.Context, arg sqlc.UpdateMagicLinkSessionParams) (sqlc.MagicLinkSession, error)
	RevokeMagicLink(ctx context.Context, token string) error

	CreateUser(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.CreateUserRow, error)
	GetUserByEmail(ctx context.Context, email string) (sqlc.GetUserByEmailRow, error)

	CreateDevice(ctx context.Context, arg sqlc.CreateDeviceParams) (sqlc.Device, error)
	GetDeviceByPubkey(ctx context.Context, pubkey string) (sqlc.Device, error)
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

type SessionStore struct {
	DeviceId string
	UserId   string
	Pubkey   string
	Email    string
}

type CooldownResult struct {
	Blocked  bool
	Cooldown *time.Time
}

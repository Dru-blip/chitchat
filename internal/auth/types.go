package auth

import (
	"chitchat/internal/db/sqlc"
	"context"

	"github.com/google/uuid"
)

type SendOtpPayload struct {
	Email  string `json:"email"`
	Pubkey string `json:"pubkey"`
}

type SendOtpResponse struct {
	Challenge string `json:"challenge"`
	Message   string `json:"message"`
	Id        string `json:"request_id"`
}

type VerifyOtpPayload struct {
	Id        string `json:"request_id"`
	Email     string `json:"email"`
	Code      string `json:"code"`
	Nonce     string `json:"n"`
	Challenge string `json:"c"`
}

type OtpSessionRespository interface {
	CreateOtpSession(ctx context.Context, arg sqlc.CreateOtpSessionParams) (sqlc.OtpSession, error)
	GetOtpSessionById(ctx context.Context, id uuid.UUID) (sqlc.GetOtpSessionByIdRow, error)
}

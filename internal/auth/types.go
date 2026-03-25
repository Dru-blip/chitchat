package auth

import (
	"chitchat/internal/db/sqlc"
	"context"

	"github.com/google/uuid"
)

type SendOtpPayload struct {
	Email  string `json:"email" validate:"required,email"`
	Pubkey string `json:"pubkey" validate:"required"`
}

type SendOtpResponse struct {
	Challenge string `json:"challenge"`
	Message   string `json:"message"`
	Id        string `json:"requestId"`
}

type VerifyOtpPayload struct {
	Id        string `json:"requestId"`
	Email     string `json:"email"`
	Code      string `json:"code"`
	Nonce     string `json:"n"`
	Challenge string `json:"c"`
}

type OtpSessionRespository interface {
	CreateOtpSession(ctx context.Context, arg sqlc.CreateOtpSessionParams) (sqlc.OtpSession, error)
	GetOtpSessionById(ctx context.Context, id uuid.UUID) (sqlc.GetOtpSessionByIdRow, error)
}

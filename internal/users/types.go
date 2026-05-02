package users

import (
	"chitchat/internal/db/sqlc"
	"context"
)

type Repository interface {
	UpdateUser(ctx context.Context, arg sqlc.UpdateUserParams) (sqlc.UpdateUserRow, error)
	OnboardUser(ctx context.Context, arg sqlc.OnboardUserParams) (sqlc.OnboardUserRow, error)
}

type OnboardUserPayload struct {
	Name      string `json:"name" validate:"required"`
	Pubkey    string `json:"pubkey" validate:"required"`
	Password  string `json:"password" validate:"required"`
	Image     string `json:"image"`
	Challenge string `json:"challenge"`
}

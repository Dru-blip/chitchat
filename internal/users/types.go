package users

import (
	"chitchat/internal/db/sqlc"
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	UpdateUser(ctx context.Context, arg sqlc.UpdateUserParams) (sqlc.UpdateUserRow, error)
	OnboardUser(ctx context.Context, arg sqlc.OnboardUserParams) (sqlc.OnboardUserRow, error)
	GetUserById(ctx context.Context, id uuid.UUID) (sqlc.GetUserByIdRow, error)
}

type OnboardUserPayload struct {
	Name      string `json:"name" validate:"required"`
	Pubkey    string `json:"pubkey" validate:"required"`
	Password  string `json:"password" validate:"required"`
	Image     string `json:"image"`
	Challenge string `json:"challenge"`
}

type UserResponse struct {
	ID         string     `json:"id"`
	Email      string     `json:"email"`
	Name       *string    `json:"name"`
	Image      *string    `json:"image"`
	CreatedAt  *time.Time `json:"createdAt"`
	Onboarding bool       `json:"onboarding"`
}

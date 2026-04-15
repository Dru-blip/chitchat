package users

import (
	"chitchat/internal/auth"
	"chitchat/internal/db/sqlc"
	"context"
)

type Service interface {
	UpadateUser(ctx context.Context) (*sqlc.User, error)
	OnboardUser(ctx context.Context, name, password, image, email string) (*sqlc.User, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) OnboardUser(ctx context.Context, name string, password string, image string, email string) (*sqlc.User, error) {
	user, err := s.repo.OnboardUser(ctx, sqlc.OnboardUserParams{
		Name:     &name,
		Image:    &image,
		Email:    email,
		Password: &password,
	})

	if err != nil {
		return nil, auth.ErrInternal
	}

	return &user, nil
}

func (s *service) UpadateUser(ctx context.Context) (*sqlc.User, error) {
	panic("unimplemented")
}

package auth

import (
	"chitchat/internal/db/sqlc"
	"chitchat/internal/utils"
	"context"
	"fmt"
	"net/netip"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	SendMagicLink(ctx context.Context, email, pubkey string, ipAddress netip.Addr, userAgent string) (*SendMagicLinkResponse, error)
	VerifyMagicLink(ctx context.Context, token string, ipAddress netip.Addr, userAgent string) (string, uuid.UUID, error)
}

type SessionInfo struct {
	SessionID string
	UserID    uuid.UUID
	Email     string
	ExpiresAt time.Time
}

type service struct {
	repo   Repository
	mailer Mailer
}

func NewService(repo Repository, mailer Mailer) Service {
	return &service{
		repo:   repo,
		mailer: mailer,
	}
}

func (s *service) SendMagicLink(ctx context.Context, email, pubkey string, ipAddress netip.Addr, userAgent string) (*SendMagicLinkResponse, error) {

	token, err := GenerateMagicLinkToken()
	if err != nil {
		return nil, ErrInternal
	}

	magic_link_session, err := s.repo.CreateMagicLinkSession(ctx, sqlc.CreateMagicLinkSessionParams{
		Email:     email,
		Pubkey:    pubkey,
		IpAddress: ipAddress,
		UserAgent: &userAgent,
		Token:     utils.SHA256(token),
		ExpiresAt: time.Now().Add(time.Minute * 15),
	})

	//TODO: get client host from env
	link := fmt.Sprintf("http://localhost:5173/verify-link?id=%s&token=%s", magic_link_session.ID, token)

	if err = s.mailer.SendMagicLink(email, link); err != nil {
		return nil, ErrInternal
	}

	return &SendMagicLinkResponse{
		Email:   email,
		Message: "successfully sent magic link",
	}, nil
}

func (s *service) VerifyMagicLink(ctx context.Context, token string, ipAddress netip.Addr, userAgent string) (string, uuid.UUID, error) {
	return "", uuid.Nil, nil
}

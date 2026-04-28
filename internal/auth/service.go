package auth

import (
	"chitchat/internal/db/sqlc"
	"chitchat/internal/utils"
	"context"
	"errors"
	"fmt"
	"net/netip"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Service interface {
	SendMagicLink(ctx context.Context, email, pubkey string, ipAddress netip.Addr, userAgent string) (*SendMagicLinkResponse, error)
	VerifyMagicLink(ctx context.Context, token string, ipAddress netip.Addr, userAgent string) (*sqlc.MagicLinkSession, error)
	GetOrCreateUser(ctx context.Context, email, pubkey string) (*sqlc.User, error)
	GetOrCreateDevice(ctx context.Context, user_id uuid.UUID, pubkey, os, user_agent string) (*sqlc.Device, bool, error)
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
	//TODO: Check for existing magic links
	token, err := utils.GenerateMagicLinkToken()
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

func (s *service) VerifyMagicLink(ctx context.Context, token string, ipAddress netip.Addr, userAgent string) (*sqlc.MagicLinkSession, error) {
	token = utils.SHA256(token)
	magic_link_session, err := s.repo.GetMagicLinkSessionByToken(ctx, token)
	session := &magic_link_session

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return session, ErrInvalidMagicLink
		}
		return session, ErrInternal
	}

	if time.Now().After(session.ExpiresAt) {
		return session, ErrMagicLinkExpired
	}

	if session.Status == sqlc.MagicLinkStatusRevoked || session.Status == sqlc.MagicLinkStatusUsed {
		return session, ErrInvalidMagicLink
	}

	_, err = s.repo.MarkMagicLinkAsUsed(ctx, token)
	if err != nil {
		return session, ErrInternal
	}

	return session, nil
}

func (s *service) GetOrCreateUser(ctx context.Context, email, pubkey string) (*sqlc.User, error) {
	row, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			createdRow, err := s.repo.CreateUser(ctx, sqlc.CreateUserParams{
				Email: email,
				Ipkey: pubkey,
			})
			if err != nil {
				return nil, err
			}
			user := sqlc.User{
				ID:         createdRow.ID,
				Email:      createdRow.Email,
				Name:       createdRow.Name,
				Image:      createdRow.Image,
				CreatedAt:  createdRow.CreatedAt,
				Onboarding: createdRow.Onboarding,
			}
			return &user, nil
		}
		return nil, err
	}
	user := sqlc.User{
		ID:         row.ID,
		Email:      row.Email,
		Name:       row.Name,
		Image:      row.Image,
		CreatedAt:  row.CreatedAt,
		Onboarding: row.Onboarding,
	}
	return &user, nil
}

func (s *service) GetOrCreateDevice(ctx context.Context, user_id uuid.UUID, pubkey, os, user_agent string) (*sqlc.Device, bool, error) {
	device, err := s.repo.GetDeviceByPubkey(ctx, pubkey)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			if device, err = s.repo.CreateDevice(ctx, sqlc.CreateDeviceParams{
				Pubkey:    pubkey,
				UserID:    user_id,
				Name:      "unknown",
				Os:        os,
				Client:    "web",
				UserAgent: &user_agent,
			}); err != nil {
				return nil, false, err
			}
			return &device, true, nil
		}
		return nil, false, err
	}
	return &device, false, nil
}

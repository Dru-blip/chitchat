package conversations

import (
	"context"
	"encoding/json"
	"fmt"

	"chitchat/internal/auth"
	"chitchat/internal/db/sqlc"

	"github.com/google/uuid"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateConversation(ctx context.Context, userID uuid.UUID, type_ string, participantEmail string) (*Conversation,
	error) {
	row, err := s.repo.CreateConversation(ctx, sqlc.CreateConversationParams{
		InitiatorID: userID,
		Type:        sqlc.ConversationTypes(type_),
		Email:       participantEmail,
	})
	fmt.Println(row, err)
	if err != nil {
		return nil, auth.ErrInternal
	}

	var participants []Participant
	if err := json.Unmarshal(row.Participants, &participants); err != nil {
		return nil, auth.ErrInternal
	}

	return &Conversation{
		ID:           row.ID,
		Type:         string(row.Type),
		Name:         row.Name,
		InitiatorID:  row.InitiatorID,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
		Participants: participants,
	}, nil
}

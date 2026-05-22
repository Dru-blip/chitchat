package conversations

import (
	"context"
	"encoding/json"
	"time"

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

	if err != nil {
		return nil, auth.ErrInternal
	}

	return toConversation(row.ID, row.Type, row.Name, row.InitiatorID, row.CreatedAt, row.UpdatedAt, row.Participants)
}

func (s *service) GetConversationsByUser(ctx context.Context, userID uuid.UUID) ([]*Conversation, error) {
	rows, err := s.repo.GetConversationsByUser(ctx, userID)
	if err != nil {
		return nil, auth.ErrInternal
	}

	var conversations []*Conversation
	for _, row := range rows {
		conv, err := toConversation(row.ID, row.Type, row.Name, row.InitiatorID, row.CreatedAt, row.UpdatedAt, row.Participants)
		if err != nil {
			return nil, err
		}
		conversations = append(conversations, conv)
	}

	return conversations, nil
}

func toConversation(id uuid.UUID, convType sqlc.ConversationTypes, name *string, initiatorID uuid.UUID, createdAt, updatedAt time.Time,
	participantsData []byte) (*Conversation, error) {
	var participants []Participant
	if err := json.Unmarshal(participantsData, &participants); err != nil {
		return nil, auth.ErrInternal
	}

	return &Conversation{
		ID:           id,
		Type:         string(convType),
		Name:         name,
		InitiatorID:  initiatorID,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		Participants: participants,
	}, nil
}

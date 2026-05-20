package conversations

import (
	"context"
	"time"

	"chitchat/internal/db/sqlc"

	"github.com/google/uuid"
)

type Conversation struct {
	ID           uuid.UUID     `json:"id"`
	Type         string        `json:"type"`
	Name         *string       `json:"name,omitempty"`
	InitiatorID  uuid.UUID     `json:"initiator_id"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	Participants []Participant `json:"participants"`
}

type Participant struct {
	ConversationID uuid.UUID  `json:"conversation_id"`
	UserID         uuid.UUID  `json:"user_id"`
	Email          string     `json:"email"`
	Image          *string    `json:"image,omitempty"`
	Name           *string    `json:"name,omitempty"`
	JoinedAt       time.Time  `json:"joined_at"`
	LeftAt         *time.Time `json:"left_at,omitempty"`
	LastRead       *time.Time `json:"last_read,omitempty"`
}

type CreateConversationPayload struct {
	Type             string `json:"type" validate:"required,oneof=dm group"`
	ParticipantEmail string `json:"participantEmail" validate:"required,email"`
}

type Repository interface {
	CreateConversation(ctx context.Context, arg sqlc.CreateConversationParams) (sqlc.CreateConversationRow, error)
}

type Service interface {
	CreateConversation(ctx context.Context, userID uuid.UUID, type_ string, participantEmail string) (*Conversation, error)
}

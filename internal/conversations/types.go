package conversations

import (
	"time"

	"github.com/google/uuid"
)

type Participant struct {
	ConversationID uuid.UUID  `json:"conversation_id"`
	UserID         uuid.UUID  `json:"user_id"`
	JoinedAt       time.Time  `json:"joined_at"`
	LeftAt         *time.Time `json:"left_at,omitempty"`
	LastRead       *time.Time `json:"last_read,omitempty"`
}

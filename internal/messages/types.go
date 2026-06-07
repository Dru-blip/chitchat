package messages

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID             uuid.UUID `json:"id"`
	ConversationID uuid.UUID `json:"conversation_id"`
	SenderID       uuid.UUID `json:"sender_id"`
	Text           string    `json:"text"`
	SentAt         time.Time `json:"sent_at"`
}

type MessageEnvelope struct {
	ConversationID uuid.UUID `json:"conversation_id"`
	MessageID      uuid.UUID `json:"id"`
	Context        string    `json:"text" validate:"required"`
	SentAt         time.Time `json:"sent_at"`
	SenderID       uuid.UUID `json:"sender_id"`
}

type SendMessageEnvelope struct {
	RecipientUserID   string `json:"recipient_user_id" validate:"required,uuid"`
	RecipientDeviceID string `json:"recipient_device_id" validate:"required,uuid"`
	IsIncoming        bool   `json:"is_incoming"`
	Context           string `json:"content" validate:"required"`
}

type SendMessagePayload struct {
	ContentType string                `json:"content_type" validate:"required"`
	Envelopes   []SendMessageEnvelope `json:"envelopes" validate:"required,min=1,dive"`
}

type Service interface {
	SendMessage(ctx context.Context, conversationID uuid.UUID, senderUserID uuid.UUID, senderDeviceID uuid.UUID, payload SendMessagePayload) (*Message, error)
}

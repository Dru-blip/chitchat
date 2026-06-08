package messages

import (
	"context"
	"fmt"
	"time"

	"chitchat/internal/db"
	"chitchat/internal/db/sqlc"

	"github.com/google/uuid"
)

type service struct {
	//TODO: should switch to Repo pattern
	// but SendMessage uses direct db,
	// so it should be divided into two objects
	// one for repo and one for transaction
	// since Im not testing the app anyways,I will leave it.
	store *db.Store
}

func NewService(store *db.Store) Service {
	return &service{store: store}
}

func (s *service) SendMessage(ctx context.Context, conversationID, senderUserID, senderDeviceID uuid.UUID, payload SendMessagePayload) (*Message, error) {
	tx, err := s.store.Db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := s.store.Queries.WithTx(tx)

	contentType := payload.ContentType

	msgRow, err := qtx.CreateMessage(ctx, sqlc.CreateMessageParams{
		ConversationID: conversationID,
		SenderUserID:   senderUserID,
		SenderDeviceID: senderDeviceID,
		ContentType:    contentType,
	})

	if err != nil {
		return nil, fmt.Errorf("create message: %w", err)
	}

	for _, envelope := range payload.Envelopes {
		recipientUserID, err := uuid.Parse(envelope.RecipientUserID)
		if err != nil {
			return nil, fmt.Errorf("invalid recipient user id: %w", err)
		}
		recipientDeviceID, err := uuid.Parse(envelope.RecipientDeviceID)
		if err != nil {
			return nil, fmt.Errorf("invalid recipient device id: %w", err)
		}

		_, err = qtx.CreateMessageEnvelope(ctx, sqlc.CreateMessageEnvelopeParams{
			MessageID:         msgRow.ID,
			RecipientUserID:   recipientUserID,
			RecipientDeviceID: recipientDeviceID,
			IsIncoming:        envelope.IsIncoming,
			Context:           envelope.Context,
		})

		if err != nil {
			return nil, fmt.Errorf("create envelope: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	return &Message{
		ID:             msgRow.ID,
		ConversationID: msgRow.ConversationID,
		SenderID:       msgRow.SenderUserID,
		SentAt:         msgRow.CreatedAt,
	}, nil
}

func (s *service) GetMessagesFromTimestamp(ctx context.Context, conversationID uuid.UUID, recipientDeviceID uuid.UUID, timestamp time.Time) ([]Message, error) {
	rows, err := s.store.Queries.GetMessagesFromTimestamp(ctx, sqlc.GetMessagesFromTimestampParams{
		ConversationID:    conversationID,
		RecipientDeviceID: recipientDeviceID,
		CreatedAt:         timestamp,
	})
	if err != nil {
		return nil, fmt.Errorf("get messages from timestamp: %w", err)
	}

	messages := make([]Message, 0, len(rows))
	for _, row := range rows {
		messages = append(messages, Message{
			ID:             row.ID,
			ConversationID: row.ConversationID,
			SenderID:       row.SenderUserID,
			SentAt:         row.CreatedAt,
			Text:           row.Context,
		})
	}

	return messages, nil
}

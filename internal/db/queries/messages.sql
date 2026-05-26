-- name: CreateMessage :one

INSERT INTO messages (conversation_id, sender_user_id, sender_device_id, sequence_id, content_type)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: CreateMessageEnvelope :one
INSERT INTO message_envelopes (message_id, recipient_user_id, recipient_device_id, is_incoming, context)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

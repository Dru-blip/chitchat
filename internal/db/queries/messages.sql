-- name: CreateMessage :one

INSERT INTO messages (conversation_id, sender_user_id, sender_device_id, sequence_id, content_type)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: CreateMessageEnvelope :one
INSERT INTO message_envelopes (message_id, recipient_user_id, recipient_device_id, is_incoming, context)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;


-- name: GetMessagesFromTimestamp :many
SELECT m.id, m.conversation_id,m.sender_user_id, me.context,m.created_at
FROM message_envelopes me
JOIN messages m ON me.message_id = m.id
WHERE m.conversation_id = $1
AND me.recipient_device_id=$2
AND m.created_at >= $3;

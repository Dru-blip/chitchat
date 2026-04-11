-- name: CreateDevice :one
INSERT INTO devices (
    pubkey, name, os, client, user_agent, user_id
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetDeviceByPubkey :one
SELECT * FROM devices
WHERE pubkey = $1
LIMIT 1;

-- name: GetDevicesByUserId :many
SELECT * FROM devices
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateDevice :one
UPDATE devices
SET name = $2, os = $3, client = $4, user_agent = $5, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

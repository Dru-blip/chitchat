-- name: CreateMagicLinkSession :one
INSERT INTO magic_link_sessions (
    token, email, pubkey, ip_address, user_agent, expires_at
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetMagicLinkSessionByToken :one
SELECT * FROM magic_link_sessions
WHERE token = $1
LIMIT 1;

-- name: GetMagicLinkSessionByEmail :many
SELECT * FROM magic_link_sessions
WHERE email = $1
ORDER BY created_at DESC;

-- name: GetPendingMagicLinkSession :one
SELECT * FROM magic_link_sessions
WHERE email = $1 
  AND status = 'pending'
  AND expires_at > CURRENT_TIMESTAMP
LIMIT 1;

-- name: UpdateMagicLinkStatus :one
UPDATE magic_link_sessions
SET status = $2
WHERE id = $1
RETURNING *;

-- name: MarkMagicLinkAsUsed :one
UPDATE magic_link_sessions
SET status = 'used'
WHERE token = $1
RETURNING *;

-- name: RevokeMagicLink :exec
UPDATE magic_link_sessions
SET status = 'revoked', updated_at = CURRENT_TIMESTAMP
WHERE token = $1;

-- name: DeleteMagicLinkSession :exec
DELETE FROM magic_link_sessions
WHERE id = $1;

-- name: DeleteMagicLinkByToken :exec
DELETE FROM magic_link_sessions
WHERE token = $1;

-- name: CleanupExpiredMagicLinks :exec
DELETE FROM magic_link_sessions
WHERE expires_at < CURRENT_TIMESTAMP;

-- name: ExpireOldMagicLinks :exec
UPDATE magic_link_sessions
SET status = 'expired'
WHERE expires_at < CURRENT_TIMESTAMP 
  AND status = 'pending';
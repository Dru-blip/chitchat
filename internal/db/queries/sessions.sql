-- name: CreateSession :one
INSERT INTO sessions (
    id, user_id, device_id, ip_address, expires_at
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetSessionById :one
SELECT * FROM sessions
WHERE id = $1 AND is_revoked = FALSE
LIMIT 1;

-- name: GetActiveSessionsByUserId :many
SELECT * FROM sessions
WHERE user_id = $1 AND is_revoked = FALSE AND expires_at > CURRENT_TIMESTAMP
ORDER BY last_active_at DESC;

-- name: UpdateSessionExpiry :exec
UPDATE sessions
SET expires_at = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND is_revoked = FALSE;

-- name: UpdateLastActiveAt :exec
UPDATE sessions
SET last_active_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: RevokeSession :exec
UPDATE sessions
SET is_revoked = TRUE, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: RevokeAllUserSessions :exec
UPDATE sessions
SET is_revoked = TRUE, updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1 AND is_revoked = FALSE;

-- name: CleanupExpiredSessions :exec
DELETE FROM sessions
WHERE expires_at < CURRENT_TIMESTAMP OR is_revoked = TRUE;

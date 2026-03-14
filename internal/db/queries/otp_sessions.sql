-- name: ListOtpSessions :many
SELECT id,
    email,
    code,
    expires_at,
    attempts,
    pubkey,
    challenge,
    created_at
FROM otp_sessions WHERE email = ?;

-- name: GetOtpSessionByEmail :one
SELECT id,
    email,
    code,
    expires_at,
    attempts,
    pubkey,
    challenge,
    created_at
FROM otp_sessions
WHERE email = ?;

-- name: GetOtpSessionByPubKey :one
SELECT id,
    email,
    code,
    expires_at,
    attempts,
    pubkey,
    challenge,
    created_at
FROM otp_sessions
WHERE pubkey = ?;


-- name: GetOtpSessionById :one
SELECT id,
    email,
    code,
    expires_at,
    attempts,
    pubkey,
    challenge,
    created_at
FROM otp_sessions
WHERE id = ?;

-- name: CreateOtpSession :one
INSERT INTO otp_sessions (id, email, code, expires_at, pubkey, challenge)
VALUES(?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateOtpSession :exec
UPDATE otp_sessions
SET expires_at = coalesce(sqlc.narg('expires_at'), expires_at),
    attempts=coalesce(sqlc.narg('attempts'),attempts)
WHERE id=sqlc.arg('id');

-- name: DeleteOtpSession :exec
DELETE FROM otp_sessions WHERE id = ?;

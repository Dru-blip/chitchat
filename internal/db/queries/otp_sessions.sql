-- name: ListOtpSessions :many
SELECT id,
    email,
    code,
    expires_at,
    attempts,
    pubkey,
    challenge,
    created_at
FROM otp_sessions WHERE email = $1;

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
WHERE email = $1;

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
WHERE pubkey = $1;


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
WHERE id = $1;

-- name: CreateOtpSession :one
INSERT INTO otp_sessions (email, code, expires_at, pubkey, challenge)
VALUES($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateOtpSession :exec
UPDATE otp_sessions
SET 
    attempts=coalesce(sqlc.narg('attempts'),attempts)
WHERE id=sqlc.arg('id');

-- name: DeleteOtpSession :exec
DELETE FROM otp_sessions WHERE id = $1;

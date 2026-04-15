-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1
LIMIT 1;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
    email, ipkey
) VALUES (
    $1, $2
)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET name = $2, image = $3, password = $4, ipkey = $5, onboarding = $6, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

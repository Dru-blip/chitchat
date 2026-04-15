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
SET name = coalesce(sqlc.narg('name'),name), 
image = coalesce(sqlc.narg('image'),image), 
password = coalesce(sqlc.narg('password'),password), 
onboarding = coalesce(sqlc.narg('onboarding'),onboarding) 
WHERE id = sqlc.arg('id')
RETURNING *;


-- name: OnboardUser :one
UPDATE users
SET name=$2,
image=coalesce(sqlc.narg('image'),image),
password=$3,
onboarding=False
WHERE email=$1 RETURNING *;

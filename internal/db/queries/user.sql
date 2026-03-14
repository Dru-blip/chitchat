-- name: ListUsers :many
SELECT id,name,email,image,created_at,ipkey FROM users ORDER BY id;

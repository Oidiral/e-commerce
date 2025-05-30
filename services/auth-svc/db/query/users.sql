-- name: CreateUser :one
INSERT INTO auth.users (email, password_hash)
VALUES ($1, $2)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM auth.users WHERE email = $1;

-- name: ListUsers :many
SELECT * FROM auth.users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

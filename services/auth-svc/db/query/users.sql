-- name: CreateUser :one
INSERT INTO users (email, password_hash)
VALUES ($1, $2)
RETURNING *;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CreateUserIfNotExists :one
INSERT INTO users (email, password_hash)
VALUES ($1, $2)
ON CONFLICT (email) DO NOTHING
RETURNING *;

-- name: GetRoleByName :one
SELECT id, name
FROM roles
WHERE name = $1;

-- name: CreateUserRole :exec
INSERT INTO user_roles (user_id, role_id)
VALUES ($1, $2)
ON CONFLICT (user_id, role_id) DO NOTHING;

-- name: GetUserByEmail :one
SELECT
    u.id,
    u.email,
    u.password_hash,
    u.status,
    u.created_at,
    u.updated_at,
    r.name AS role_name
FROM users u
JOIN user_roles ur ON ur.user_id = u.id
JOIN roles r ON r.id = ur.role_id
WHERE u.email = $1;

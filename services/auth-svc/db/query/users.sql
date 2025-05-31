-- name: CreateUser :one
INSERT INTO auth.users (email, password_hash)
VALUES ($1, $2)
RETURNING *;



-- name: ListUsers :many
SELECT * FROM auth.users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CreateUserIfNotExists :one
INSERT INTO auth.users (email, password_hash)
VALUES ($1, $2)
ON CONFLICT (email) DO NOTHING
RETURNING *;

-- name: GetRoleByName :one
SELECT id, name
FROM auth.roles
WHERE name = $1;


-- name: CreateUserRole :exec
INSERT INTO auth.user_roles (user_id, role_id)
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
FROM auth.users u
JOIN auth.user_roles ur ON ur.user_id = u.id
JOIN auth.roles r ON r.id = ur.role_id
WHERE u.email = $1;


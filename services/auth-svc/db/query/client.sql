-- name: GetById :one
SELECT c.id,
       c.secret_hash,
       c.roles,
       c.status,
       c.created_at
FROM clients c WHERE id = $1;
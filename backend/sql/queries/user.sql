-- name: UpsertUser :exec
INSERT INTO users (id, email, name, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (id) DO UPDATE SET
    email      = EXCLUDED.email,
    name       = EXCLUDED.name,
    updated_at = EXCLUDED.updated_at;

-- name: FindUserByID :one
SELECT id, email, name, created_at, updated_at
FROM users
WHERE id = $1;

-- name: FindUserByEmail :one
SELECT id, email, name, created_at, updated_at
FROM users
WHERE email = $1;

-- name: DeleteUser :execrows
DELETE FROM users
WHERE id = $1;

-- name: UpsertCompany :exec
INSERT INTO companies (id, user_id, name, memo, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (id) DO UPDATE SET
    name       = EXCLUDED.name,
    memo       = EXCLUDED.memo,
    updated_at = EXCLUDED.updated_at;

-- name: FindCompanyByID :one
SELECT id, user_id, name, memo, created_at, updated_at
FROM companies
WHERE user_id = $1 AND id = $2;

-- name: ListCompaniesByUserID :many
SELECT id, user_id, name, memo, created_at, updated_at
FROM companies
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: DeleteCompany :execrows
DELETE FROM companies
WHERE user_id = $1 AND id = $2;

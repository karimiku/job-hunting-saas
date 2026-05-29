-- name: CreateCompanyAlias :exec
INSERT INTO company_aliases (id, user_id, company_id, alias, created_at)
VALUES ($1, $2, $3, $4, $5);

-- name: FindCompanyAliasByID :one
SELECT id, user_id, company_id, alias, created_at
FROM company_aliases
WHERE user_id = $1 AND id = $2;

-- name: ListCompanyAliasesByCompanyID :many
SELECT id, user_id, company_id, alias, created_at
FROM company_aliases
WHERE user_id = $1 AND company_id = $2
ORDER BY created_at DESC;

-- name: ListCompanyAliasesByUserID :many
SELECT id, user_id, company_id, alias, created_at
FROM company_aliases
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: DeleteCompanyAlias :execrows
DELETE FROM company_aliases
WHERE user_id = $1 AND id = $2;

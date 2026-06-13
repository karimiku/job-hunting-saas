-- name: CreateAIAccessToken :exec
INSERT INTO ai_access_tokens (
    id,
    user_id,
    name,
    token_hash,
    token_preview,
    last_used_at,
    revoked_at,
    created_at,
    updated_at
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9
);

-- name: FindActiveAIAccessTokenByHash :one
SELECT
    id,
    user_id,
    name,
    token_hash,
    token_preview,
    last_used_at,
    revoked_at,
    created_at,
    updated_at
FROM ai_access_tokens
WHERE token_hash = $1
  AND revoked_at IS NULL;

-- name: ListAIAccessTokensByUserID :many
SELECT
    id,
    user_id,
    name,
    token_hash,
    token_preview,
    last_used_at,
    revoked_at,
    created_at,
    updated_at
FROM ai_access_tokens
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: TouchAIAccessTokenLastUsed :execrows
UPDATE ai_access_tokens
SET
    last_used_at = $2,
    updated_at = $2
WHERE id = $1
  AND revoked_at IS NULL;

-- name: RevokeAIAccessToken :execrows
UPDATE ai_access_tokens
SET
    revoked_at = $3,
    updated_at = $3
WHERE user_id = $1
  AND id = $2
  AND revoked_at IS NULL;

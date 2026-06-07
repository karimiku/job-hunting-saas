-- name: CreateESMemo :one
INSERT INTO es_memos (id, user_id, entry_id, category, title, content, source, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, user_id, entry_id, category, title, content, source, created_at, updated_at;

-- name: ListESMemosByUserID :many
SELECT id, user_id, entry_id, category, title, content, source, created_at, updated_at
FROM es_memos
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2;

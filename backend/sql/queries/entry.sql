-- name: UpsertEntry :exec
INSERT INTO entries (id, user_id, company_id, route, source, status, stage_kind, stage_label, memo, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
ON CONFLICT (id) DO UPDATE SET
    route       = EXCLUDED.route,
    source      = EXCLUDED.source,
    status      = EXCLUDED.status,
    stage_kind  = EXCLUDED.stage_kind,
    stage_label = EXCLUDED.stage_label,
    memo        = EXCLUDED.memo,
    updated_at  = EXCLUDED.updated_at;

-- name: FindEntryByID :one
SELECT id, user_id, company_id, route, source, status, stage_kind, stage_label, memo, created_at, updated_at
FROM entries
WHERE user_id = $1 AND id = $2;

-- name: ListEntriesByUserID :many
SELECT id, user_id, company_id, route, source, status, stage_kind, stage_label, memo, created_at, updated_at
FROM entries
WHERE user_id = $1
  AND (sqlc.narg('status')::entry_status IS NULL OR status = sqlc.narg('status'))
  AND (sqlc.narg('stage_kind')::stage_kind IS NULL OR stage_kind = sqlc.narg('stage_kind'))
  AND (sqlc.narg('source')::text IS NULL OR source = sqlc.narg('source'))
ORDER BY updated_at DESC;

-- name: DeleteEntry :execrows
DELETE FROM entries
WHERE user_id = $1 AND id = $2;

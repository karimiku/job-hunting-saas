-- name: CreateStageHistory :exec
INSERT INTO stage_histories (id, entry_id, stage_kind, stage_label, note, created_at)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: ListStageHistoriesByEntryID :many
SELECT id, entry_id, stage_kind, stage_label, note, created_at
FROM stage_histories
WHERE entry_id = $1
ORDER BY created_at;

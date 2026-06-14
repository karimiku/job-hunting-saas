-- name: CreateStageHistory :exec
INSERT INTO stage_histories (id, entry_id, stage_kind, stage_label, note, created_at)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: ListStageHistoriesByEntryID :many
SELECT sh.id, sh.entry_id, sh.stage_kind, sh.stage_label, sh.note, sh.created_at
FROM stage_histories sh
JOIN entries e ON e.id = sh.entry_id
WHERE e.user_id = $1 AND sh.entry_id = $2
ORDER BY sh.created_at;

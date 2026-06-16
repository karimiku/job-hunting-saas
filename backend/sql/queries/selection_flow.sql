-- name: UpsertSelectionFlow :one
INSERT INTO selection_flows (
    id,
    entry_id,
    source,
    current_stage_position,
    confidence,
    inbox_clip_id,
    created_at,
    updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (entry_id) DO UPDATE SET
    source = EXCLUDED.source,
    current_stage_position = EXCLUDED.current_stage_position,
    confidence = EXCLUDED.confidence,
    inbox_clip_id = EXCLUDED.inbox_clip_id,
    updated_at = EXCLUDED.updated_at
RETURNING id, entry_id, source, current_stage_position, confidence, inbox_clip_id, created_at, updated_at;

-- name: DeleteSelectionStagesByFlowID :exec
DELETE FROM selection_stages
WHERE flow_id = $1;

-- name: CreateSelectionStage :exec
INSERT INTO selection_stages (id, flow_id, position, stage_kind, stage_label, evidence_text, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: FindSelectionFlowByEntryID :one
SELECT sf.id, sf.entry_id, sf.source, sf.current_stage_position, sf.confidence, sf.inbox_clip_id, sf.created_at, sf.updated_at
FROM selection_flows sf
JOIN entries e ON e.id = sf.entry_id
WHERE e.user_id = $1 AND sf.entry_id = $2;

-- name: ListSelectionStagesByFlowID :many
SELECT id, flow_id, position, stage_kind, stage_label, evidence_text, created_at
FROM selection_stages
WHERE flow_id = $1
ORDER BY position ASC;

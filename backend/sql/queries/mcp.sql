-- name: MCPListEntries :many
SELECT
    e.id,
    e.company_id,
    c.name AS company_name,
    e.route,
    e.source,
    e.source_url,
    e.status,
    e.stage_kind,
    e.stage_label,
    e.memo,
    e.created_at,
    e.updated_at
FROM entries e
JOIN companies c ON c.id = e.company_id
WHERE e.user_id = $1
ORDER BY e.updated_at DESC;

-- name: MCPGetEntryContext :one
SELECT
    e.id,
    e.company_id,
    c.name AS company_name,
    e.route,
    e.source,
    e.source_url,
    e.status,
    e.stage_kind,
    e.stage_label,
    e.memo,
    e.created_at,
    e.updated_at
FROM entries e
JOIN companies c ON c.id = e.company_id
WHERE e.user_id = $1 AND e.id = $2;

-- name: MCPListOpenTasks :many
SELECT
    t.id,
    t.entry_id,
    c.name AS company_name,
    t.title,
    t.task_type,
    t.due_date,
    t.status,
    t.notify,
    t.memo,
    t.created_at,
    t.updated_at
FROM tasks t
JOIN entries e ON e.id = t.entry_id
JOIN companies c ON c.id = e.company_id
WHERE e.user_id = $1 AND t.status = 'todo'
ORDER BY t.due_date ASC NULLS LAST, t.created_at ASC;

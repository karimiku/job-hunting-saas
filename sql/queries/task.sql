-- name: UpsertTask :exec
INSERT INTO tasks (id, entry_id, title, task_type, due_date, status, notify, memo, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
ON CONFLICT (id) DO UPDATE SET
    title      = EXCLUDED.title,
    task_type  = EXCLUDED.task_type,
    due_date   = EXCLUDED.due_date,
    status     = EXCLUDED.status,
    notify     = EXCLUDED.notify,
    memo       = EXCLUDED.memo,
    updated_at = EXCLUDED.updated_at;

-- name: FindTaskByID :one
SELECT t.id, t.entry_id, t.title, t.task_type, t.due_date, t.status, t.notify, t.memo, t.created_at, t.updated_at
FROM tasks t
JOIN entries e ON t.entry_id = e.id
WHERE e.user_id = $1 AND t.id = $2;

-- name: ListTasksByEntryID :many
SELECT t.id, t.entry_id, t.title, t.task_type, t.due_date, t.status, t.notify, t.memo, t.created_at, t.updated_at
FROM tasks t
JOIN entries e ON t.entry_id = e.id
WHERE e.user_id = $1 AND t.entry_id = $2
ORDER BY t.created_at;

-- name: ListTasksByUserIDWithDueBefore :many
SELECT t.id, t.entry_id, t.title, t.task_type, t.due_date, t.status, t.notify, t.memo, t.created_at, t.updated_at
FROM tasks t
JOIN entries e ON t.entry_id = e.id
WHERE e.user_id = $1 AND t.status = 'todo' AND t.due_date < $2
ORDER BY t.due_date;

-- name: DeleteTask :execrows
DELETE FROM tasks
USING entries e
WHERE tasks.entry_id = e.id AND e.user_id = $1 AND tasks.id = $2;

-- name: CreateInboxClip :exec
INSERT INTO inbox_clips (id, user_id, url, title, source, guess, captured_at)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: FindInboxClipByID :one
SELECT id, user_id, url, title, source, guess, captured_at
FROM inbox_clips
WHERE user_id = $1 AND id = $2;

-- name: ListInboxClipsByUserID :many
SELECT id, user_id, url, title, source, guess, captured_at
FROM inbox_clips
WHERE user_id = $1
ORDER BY captured_at DESC;

-- name: DeleteInboxClip :execrows
DELETE FROM inbox_clips
WHERE user_id = $1 AND id = $2;

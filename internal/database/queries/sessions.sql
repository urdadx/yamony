-- name: CreateSession :one
INSERT INTO sessions (user_id, session_token, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetSessionByToken :one
SELECT * FROM sessions
WHERE session_token = $1;

-- name: GetSessionsByUserID :many
SELECT * FROM sessions
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = $1;

-- name: DeleteExpiredSessions :exec
DELETE FROM sessions
WHERE expires_at < NOW();

-- name: DeleteUserSessions :exec
DELETE FROM sessions
WHERE user_id = $1;


-- name: UpdateSessionWithActivePage :one
UPDATE sessions
SET active_page_id = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

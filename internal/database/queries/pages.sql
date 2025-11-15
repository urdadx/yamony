-- name: CreatePage :one
INSERT INTO pages (
    user_id,
    name,
    handle,
    banner_image,
    image,
    bio,
    is_active
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetPageByID :one
SELECT * FROM pages
WHERE id = $1;

-- name: GetPageByHandle :one
SELECT * FROM pages
WHERE handle = $1;

-- name: GetPagesByUserID :many
SELECT * FROM pages
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetActivePagesByUserID :many
SELECT * FROM pages
WHERE user_id = $1 AND is_active = true
ORDER BY created_at DESC;

-- name: UpdatePage :one
UPDATE pages
SET
    name = COALESCE(sqlc.narg('name'), name),
    handle = COALESCE(sqlc.narg('handle'), handle),
    banner_image = COALESCE(sqlc.narg('banner_image'), banner_image),
    image = COALESCE(sqlc.narg('image'), image),
    bio = COALESCE(sqlc.narg('bio'), bio),
    is_active = COALESCE(sqlc.narg('is_active'), is_active),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeletePage :exec
DELETE FROM pages
WHERE id = $1;

-- name: CheckHandleExists :one
SELECT EXISTS(
    SELECT 1 FROM pages
    WHERE handle = $1 AND id != COALESCE($2, 0)
);

-- name: GetAllPages :many
SELECT * FROM pages
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountUserPages :one
SELECT COUNT(*) FROM pages
WHERE user_id = $1;

-- name: GetUserMostRecentPage :one
SELECT * FROM pages
WHERE user_id = $1
ORDER BY updated_at DESC
LIMIT 1;

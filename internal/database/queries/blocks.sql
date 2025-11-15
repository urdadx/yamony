-- name: CreateBlock :one
INSERT INTO blocks (
    page_id,
    user_id,
    block_order,
    is_active,
    title,
    description,
    layout,
    block_type,
    variant,
    properties
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING *;

-- name: GetBlockByID :one
SELECT * FROM blocks
WHERE id = $1;

-- name: GetBlocksByPageID :many
SELECT * FROM blocks
WHERE page_id = $1
ORDER BY block_order ASC NULLS LAST, created_at ASC;

-- name: GetActiveBlocksByPageID :many
SELECT * FROM blocks
WHERE page_id = $1 AND is_active = true
ORDER BY block_order ASC NULLS LAST, created_at ASC;

-- name: GetBlocksByUserID :many
SELECT * FROM blocks
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetBlocksByPageIDAndType :many
SELECT * FROM blocks
WHERE page_id = $1 AND block_type = $2
ORDER BY block_order ASC NULLS LAST, created_at ASC;

-- name: UpdateBlock :one
UPDATE blocks
SET
    block_order = COALESCE(sqlc.narg('order'), block_order),
    is_active = COALESCE(sqlc.narg('is_active'), is_active),
    title = COALESCE(sqlc.narg('title'), title),
    description = COALESCE(sqlc.narg('description'), description),
    layout = COALESCE(sqlc.narg('layout'), layout),
    block_type = COALESCE(sqlc.narg('block_type'), block_type),
    variant = COALESCE(sqlc.narg('variant'), variant),
    properties = COALESCE(sqlc.narg('properties'), properties),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: UpdateBlockOrder :exec
UPDATE blocks
SET block_order = $2, updated_at = NOW()
WHERE id = $1;

-- name: DeleteBlock :exec
DELETE FROM blocks
WHERE id = $1;

-- name: DeleteBlocksByPageID :exec
DELETE FROM blocks
WHERE page_id = $1;

-- name: CountBlocksByPageID :one
SELECT COUNT(*) FROM blocks
WHERE page_id = $1;

-- name: CountActiveBlocksByPageID :one
SELECT COUNT(*) FROM blocks
WHERE page_id = $1 AND is_active = true;


-- name: ReorderBlocks :exec
UPDATE blocks
SET block_order = block_order + $3, updated_at = NOW()
WHERE page_id = $1 AND block_order >= $2;


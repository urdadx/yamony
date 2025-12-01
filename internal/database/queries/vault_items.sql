-- name: CreateVaultItem :one
INSERT INTO vault_items (
    id,
    vault_id,
    item_type,
    encrypted_blob,
    iv,
    tag,
    meta,
    version
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetVaultItemByID :one
SELECT * FROM vault_items
WHERE id = $1;

-- name: GetVaultItemsByVaultID :many
SELECT * FROM vault_items
WHERE vault_id = $1
ORDER BY created_at DESC;

-- name: GetVaultItemsByVaultIDAndType :many
SELECT * FROM vault_items
WHERE vault_id = $1 AND item_type = $2
ORDER BY created_at DESC;

-- name: UpdateVaultItem :one
UPDATE vault_items
SET 
    encrypted_blob = $2,
    iv = $3,
    tag = $4,
    meta = $5,
    version = $6,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteVaultItem :exec
DELETE FROM vault_items
WHERE id = $1;

-- name: GetVaultItemsByIDs :many
SELECT * FROM vault_items
WHERE id = ANY($1::uuid[]);

-- name: CountVaultItems :one
SELECT COUNT(*) FROM vault_items
WHERE vault_id = $1;

-- name: SearchVaultItemsByMeta :many
SELECT * FROM vault_items
WHERE vault_id = $1 
  AND meta @> $2::jsonb
ORDER BY created_at DESC
LIMIT $3;

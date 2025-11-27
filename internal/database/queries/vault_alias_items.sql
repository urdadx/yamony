-- name: CreateAliasItem :one
INSERT INTO vault_alias_items (
    vault_id,
    user_id,
    title,
    alias_prefix,
    alias_suffix,
    forwards_to,
    note,
    is_favorite
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetAliasItemByID :one
SELECT * FROM vault_alias_items
WHERE id = $1 AND user_id = $2;

-- name: GetVaultAliasItems :many
SELECT * FROM vault_alias_items
WHERE vault_id = $1 AND user_id = $2
ORDER BY is_favorite DESC, updated_at DESC;

-- name: GetUserAliasItems :many
SELECT * FROM vault_alias_items
WHERE user_id = $1
ORDER BY is_favorite DESC, updated_at DESC;

-- name: UpdateAliasItem :one
UPDATE vault_alias_items
SET 
    title = $2,
    alias_prefix = $3,
    alias_suffix = $4,
    forwards_to = $5,
    note = $6,
    is_favorite = $7,
    updated_at = NOW()
WHERE id = $1 AND user_id = $8
RETURNING *;

-- name: UpdateAliasItemLastUsed :exec
UPDATE vault_alias_items
SET last_used_at = NOW()
WHERE id = $1 AND user_id = $2;

-- name: DeleteAliasItem :exec
DELETE FROM vault_alias_items
WHERE id = $1 AND user_id = $2;

-- name: ToggleAliasItemFavorite :one
UPDATE vault_alias_items
SET 
    is_favorite = NOT is_favorite,
    updated_at = NOW()
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: SearchAliasItems :many
SELECT * FROM vault_alias_items
WHERE user_id = $1 AND (
    title ILIKE '%' || $2 || '%' OR
    alias_prefix ILIKE '%' || $2 || '%' OR
    alias_suffix ILIKE '%' || $2 || '%' OR
    forwards_to ILIKE '%' || $2 || '%' OR
    note ILIKE '%' || $2 || '%'
)
ORDER BY is_favorite DESC, updated_at DESC;

-- Alias Attachments queries
-- name: AddAliasAttachment :one
INSERT INTO vault_alias_attachments (
    alias_item_id,
    filename,
    file_path,
    file_size,
    mime_type
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetAliasAttachments :many
SELECT * FROM vault_alias_attachments
WHERE alias_item_id = $1
ORDER BY created_at ASC;

-- name: GetAliasAttachmentByID :one
SELECT * FROM vault_alias_attachments
WHERE id = $1 AND alias_item_id = $2;

-- name: DeleteAliasAttachment :exec
DELETE FROM vault_alias_attachments
WHERE id = $1 AND alias_item_id = $2;

-- name: DeleteAllAliasAttachments :exec
DELETE FROM vault_alias_attachments
WHERE alias_item_id = $1;

-- name: CreateLoginItem :one
INSERT INTO vault_login_items (
    vault_id,
    user_id,
    title,
    username,
    password_encrypted,
    totp_secret_encrypted,
    note,
    is_favorite
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetLoginItemByID :one
SELECT * FROM vault_login_items
WHERE id = $1 AND user_id = $2;

-- name: GetVaultLoginItems :many
SELECT * FROM vault_login_items
WHERE vault_id = $1 AND user_id = $2
ORDER BY is_favorite DESC, updated_at DESC;

-- name: GetUserLoginItems :many
SELECT * FROM vault_login_items
WHERE user_id = $1
ORDER BY is_favorite DESC, updated_at DESC;

-- name: UpdateLoginItem :one
UPDATE vault_login_items
SET 
    title = $2,
    username = $3,
    password_encrypted = $4,
    totp_secret_encrypted = $5,
    note = $6,
    is_favorite = $7,
    updated_at = NOW()
WHERE id = $1 AND user_id = $8
RETURNING *;

-- name: UpdateLoginItemLastUsed :exec
UPDATE vault_login_items
SET last_used_at = NOW()
WHERE id = $1 AND user_id = $2;

-- name: DeleteLoginItem :exec
DELETE FROM vault_login_items
WHERE id = $1 AND user_id = $2;

-- name: ToggleLoginItemFavorite :one
UPDATE vault_login_items
SET 
    is_favorite = NOT is_favorite,
    updated_at = NOW()
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: SearchLoginItems :many
SELECT * FROM vault_login_items
WHERE user_id = $1 AND (
    title ILIKE '%' || $2 || '%' OR
    username ILIKE '%' || $2 || '%' OR
    note ILIKE '%' || $2 || '%'
)
ORDER BY is_favorite DESC, updated_at DESC;

-- Login Websites queries
-- name: AddLoginWebsite :one
INSERT INTO vault_login_websites (
    login_item_id,
    url
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetLoginWebsites :many
SELECT * FROM vault_login_websites
WHERE login_item_id = $1
ORDER BY created_at ASC;

-- name: DeleteLoginWebsite :exec
DELETE FROM vault_login_websites
WHERE id = $1 AND login_item_id = $2;

-- name: DeleteAllLoginWebsites :exec
DELETE FROM vault_login_websites
WHERE login_item_id = $1;

-- Login Attachments queries
-- name: AddLoginAttachment :one
INSERT INTO vault_login_attachments (
    login_item_id,
    filename,
    file_path,
    file_size,
    mime_type
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetLoginAttachments :many
SELECT * FROM vault_login_attachments
WHERE login_item_id = $1
ORDER BY created_at ASC;

-- name: GetLoginAttachmentByID :one
SELECT * FROM vault_login_attachments
WHERE id = $1 AND login_item_id = $2;

-- name: DeleteLoginAttachment :exec
DELETE FROM vault_login_attachments
WHERE id = $1 AND login_item_id = $2;

-- name: DeleteAllLoginAttachments :exec
DELETE FROM vault_login_attachments
WHERE login_item_id = $1;

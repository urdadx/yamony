-- name: CreateNoteItem :one
INSERT INTO vault_note_items (
    vault_id,
    user_id,
    title,
    note,
    is_favorite
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetNoteItemByID :one
SELECT * FROM vault_note_items
WHERE id = $1 AND user_id = $2;

-- name: GetVaultNoteItems :many
SELECT * FROM vault_note_items
WHERE vault_id = $1 AND user_id = $2
ORDER BY is_favorite DESC, updated_at DESC;

-- name: GetUserNoteItems :many
SELECT * FROM vault_note_items
WHERE user_id = $1
ORDER BY is_favorite DESC, updated_at DESC;

-- name: UpdateNoteItem :one
UPDATE vault_note_items
SET 
    title = $2,
    note = $3,
    is_favorite = $4,
    updated_at = NOW()
WHERE id = $1 AND user_id = $5
RETURNING *;

-- name: UpdateNoteItemLastUsed :exec
UPDATE vault_note_items
SET last_used_at = NOW()
WHERE id = $1 AND user_id = $2;

-- name: DeleteNoteItem :exec
DELETE FROM vault_note_items
WHERE id = $1 AND user_id = $2;

-- name: ToggleNoteItemFavorite :one
UPDATE vault_note_items
SET 
    is_favorite = NOT is_favorite,
    updated_at = NOW()
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: SearchNoteItems :many
SELECT * FROM vault_note_items
WHERE user_id = $1 AND (
    title ILIKE '%' || $2 || '%' OR
    note ILIKE '%' || $2 || '%'
)
ORDER BY is_favorite DESC, updated_at DESC;

-- Note Attachments queries
-- name: AddNoteAttachment :one
INSERT INTO vault_note_attachments (
    note_item_id,
    filename,
    file_path,
    file_size,
    mime_type
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetNoteAttachments :many
SELECT * FROM vault_note_attachments
WHERE note_item_id = $1
ORDER BY created_at ASC;

-- name: GetNoteAttachmentByID :one
SELECT * FROM vault_note_attachments
WHERE id = $1 AND note_item_id = $2;

-- name: DeleteNoteAttachment :exec
DELETE FROM vault_note_attachments
WHERE id = $1 AND note_item_id = $2;

-- name: DeleteAllNoteAttachments :exec
DELETE FROM vault_note_attachments
WHERE note_item_id = $1;

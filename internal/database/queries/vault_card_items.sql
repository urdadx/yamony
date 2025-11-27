-- name: CreateCardItem :one
INSERT INTO vault_card_items (
    vault_id,
    user_id,
    name,
    card_number_encrypted,
    expiration_date,
    security_code_encrypted,
    pin_encrypted,
    note,
    is_favorite
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetCardItemByID :one
SELECT * FROM vault_card_items
WHERE id = $1 AND user_id = $2;

-- name: GetVaultCardItems :many
SELECT * FROM vault_card_items
WHERE vault_id = $1 AND user_id = $2
ORDER BY is_favorite DESC, updated_at DESC;

-- name: GetUserCardItems :many
SELECT * FROM vault_card_items
WHERE user_id = $1
ORDER BY is_favorite DESC, updated_at DESC;

-- name: UpdateCardItem :one
UPDATE vault_card_items
SET 
    name = $2,
    card_number_encrypted = $3,
    expiration_date = $4,
    security_code_encrypted = $5,
    pin_encrypted = $6,
    note = $7,
    is_favorite = $8,
    updated_at = NOW()
WHERE id = $1 AND user_id = $9
RETURNING *;

-- name: UpdateCardItemLastUsed :exec
UPDATE vault_card_items
SET last_used_at = NOW()
WHERE id = $1 AND user_id = $2;

-- name: DeleteCardItem :exec
DELETE FROM vault_card_items
WHERE id = $1 AND user_id = $2;

-- name: ToggleCardItemFavorite :one
UPDATE vault_card_items
SET 
    is_favorite = NOT is_favorite,
    updated_at = NOW()
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: SearchCardItems :many
SELECT * FROM vault_card_items
WHERE user_id = $1 AND (
    name ILIKE '%' || $2 || '%' OR
    note ILIKE '%' || $2 || '%'
)
ORDER BY is_favorite DESC, updated_at DESC;

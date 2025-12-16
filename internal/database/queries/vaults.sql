-- name: CreateVault :one
INSERT INTO vaults (
    user_id,
    name,
    description,
    icon,
    theme,
    is_favorite
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetVaultByID :one
SELECT * FROM vaults
WHERE id = $1;

-- name: GetVaultByIDAndUserID :one
SELECT * FROM vaults
WHERE id = $1 AND user_id = $2;

-- name: GetUserVaults :many
SELECT 
    v.*,
    COALESCE(COUNT(vi.id), 0)::int AS item_count
FROM vaults v
LEFT JOIN vault_items vi ON v.id = vi.vault_id
WHERE v.user_id = $1
GROUP BY v.id
ORDER BY v.is_favorite DESC, v.updated_at DESC;

-- name: UpdateVault :one
UPDATE vaults
SET 
    name = $2,
    description = $3,
    icon = $4,
    theme = $5,
    is_favorite = $6,
    updated_at = NOW()
WHERE id = $1 AND user_id = $7
RETURNING *;

-- name: DeleteVault :exec
DELETE FROM vaults
WHERE id = $1 AND user_id = $2;

-- name: ToggleVaultFavorite :one
UPDATE vaults
SET 
    is_favorite = NOT is_favorite,
    updated_at = NOW()
WHERE id = $1 AND user_id = $2
RETURNING *;

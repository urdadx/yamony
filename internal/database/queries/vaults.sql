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
WHERE id = $1 AND user_id = $2;

-- name: GetUserVaults :many
SELECT * FROM vaults
WHERE user_id = $1
ORDER BY is_favorite DESC, updated_at DESC;

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

-- name: CreateVaultKey :one
INSERT INTO vault_keys (
    vault_id,
    wrapped_vek,
    wrap_iv,
    wrap_tag,
    kdf_salt,
    kdf_params,
    version,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, NOW(), NOW()
)
RETURNING *;

-- name: GetVaultKeyByVaultID :one
SELECT * FROM vault_keys
WHERE vault_id = $1
ORDER BY version DESC
LIMIT 1;

-- name: GetVaultKeyByVaultIDAndVersion :one
SELECT * FROM vault_keys
WHERE vault_id = $1 AND version = $2
LIMIT 1;

-- name: GetAllVaultKeyVersions :many
SELECT * FROM vault_keys
WHERE vault_id = $1
ORDER BY version DESC;

-- name: UpdateVaultKey :one
UPDATE vault_keys
SET 
    wrapped_vek = $2,
    wrap_iv = $3,
    wrap_tag = $4,
    kdf_salt = $5,
    kdf_params = $6,
    version = $7,
    updated_at = NOW()
WHERE vault_id = $1
RETURNING *;

-- name: DeleteVaultKeys :exec
DELETE FROM vault_keys
WHERE vault_id = $1;

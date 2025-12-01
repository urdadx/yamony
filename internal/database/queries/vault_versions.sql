-- name: CreateVaultVersion :one
INSERT INTO vault_versions (
    vault_id,
    object_key,
    mac,
    created_by_device
) VALUES (
    $1, $2, $3, $4
)
RETURNING id, vault_id, object_key, mac, created_by_device, created_at;

-- name: GetVaultVersionByID :one
SELECT id, vault_id, object_key, mac, created_by_device, created_at
FROM vault_versions
WHERE id = $1;

-- name: GetLatestVaultVersion :one
SELECT id, vault_id, object_key, mac, created_by_device, created_at
FROM vault_versions
WHERE vault_id = $1
ORDER BY id DESC
LIMIT 1;

-- name: GetVaultVersionsByVaultID :many
SELECT id, vault_id, object_key, mac, created_by_device, created_at
FROM vault_versions
WHERE vault_id = $1
ORDER BY id DESC
LIMIT $2;

-- name: GetVaultVersionByIDAndVault :one
SELECT id, vault_id, object_key, mac, created_by_device, created_at
FROM vault_versions
WHERE vault_id = $1 AND id = $2;

-- name: GetVaultVersionsSinceID :many
SELECT id, vault_id, object_key, mac, created_by_device, created_at
FROM vault_versions
WHERE vault_id = $1 AND id > $2
ORDER BY id ASC;

-- name: CountVaultVersions :one
SELECT COUNT(*) FROM vault_versions
WHERE vault_id = $1;

-- name: DeleteOldVaultVersions :exec
DELETE FROM vault_versions v1
WHERE v1.vault_id = $1 
AND v1.id < (
    SELECT MAX(v2.id) - $2 
    FROM vault_versions v2
    WHERE v2.vault_id = $1
);

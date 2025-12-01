-- name: CreateDevice :one
INSERT INTO devices (
    id,
    user_id,
    device_label,
    x25519_public,
    ed25519_public,
    created_at,
    last_seen
) VALUES (
    $1, $2, $3, $4, $5, NOW(), NOW()
)
RETURNING *;

-- name: GetDeviceByID :one
SELECT * FROM devices
WHERE id = $1 AND revoked_at IS NULL
LIMIT 1;

-- name: GetDevicesByUserID :many
SELECT * FROM devices
WHERE user_id = $1 AND revoked_at IS NULL
ORDER BY created_at DESC;

-- name: GetAllDevicesByUserID :many
SELECT * FROM devices
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateDeviceLastSeen :exec
UPDATE devices
SET last_seen = NOW()
WHERE id = $1;

-- name: RevokeDevice :exec
UPDATE devices
SET revoked_at = NOW()
WHERE id = $1;

-- name: GetDevicePublicKeys :one
SELECT x25519_public, ed25519_public
FROM devices
WHERE id = $1 AND revoked_at IS NULL
LIMIT 1;

-- name: GetUserDevicePublicKeys :many
SELECT id, device_label, x25519_public, ed25519_public, created_at
FROM devices
WHERE user_id = $1 AND revoked_at IS NULL
ORDER BY created_at DESC;

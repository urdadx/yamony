-- name: CreateSharingRecord :one
INSERT INTO sharing_records (
    vault_id,
    item_id,
    sender_user_id,
    recipient_user_id,
    wrapped_key,
    wrap_iv,
    wrap_tag,
    status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING id, vault_id, item_id, sender_user_id, recipient_user_id, wrapped_key, wrap_iv, wrap_tag, status, created_at, accepted_at;

-- name: GetSharingRecordByID :one
SELECT id, vault_id, item_id, sender_user_id, recipient_user_id, wrapped_key, wrap_iv, wrap_tag, status, created_at, accepted_at
FROM sharing_records
WHERE id = $1;

-- name: GetSharingRecordsByVaultID :many
SELECT id, vault_id, item_id, sender_user_id, recipient_user_id, wrapped_key, wrap_iv, wrap_tag, status, created_at, accepted_at
FROM sharing_records
WHERE vault_id = $1 AND status = 'accepted'
ORDER BY created_at DESC;

-- name: GetSharingRecordsByRecipientID :many
SELECT id, vault_id, item_id, sender_user_id, recipient_user_id, wrapped_key, wrap_iv, wrap_tag, status, created_at, accepted_at
FROM sharing_records
WHERE recipient_user_id = $1
ORDER BY created_at DESC;

-- name: GetPendingSharingRecordsByRecipientID :many
SELECT id, vault_id, item_id, sender_user_id, recipient_user_id, wrapped_key, wrap_iv, wrap_tag, status, created_at, accepted_at
FROM sharing_records
WHERE recipient_user_id = $1 AND status = 'pending'
ORDER BY created_at DESC;

-- name: GetSharingRecordByVaultAndRecipient :one
SELECT id, vault_id, item_id, sender_user_id, recipient_user_id, wrapped_key, wrap_iv, wrap_tag, status, created_at, accepted_at
FROM sharing_records
WHERE vault_id = $1 AND recipient_user_id = $2 AND status = 'accepted';

-- name: AcceptSharingRecord :one
UPDATE sharing_records
SET status = 'accepted', accepted_at = NOW()
WHERE id = $1 AND status = 'pending'
RETURNING id, vault_id, item_id, sender_user_id, recipient_user_id, wrapped_key, wrap_iv, wrap_tag, status, created_at, accepted_at;

-- name: RejectSharingRecord :one
UPDATE sharing_records
SET status = 'rejected'
WHERE id = $1 AND status = 'pending'
RETURNING id, vault_id, item_id, sender_user_id, recipient_user_id, wrapped_key, wrap_iv, wrap_tag, status, created_at, accepted_at;

-- name: RevokeSharingRecord :exec
UPDATE sharing_records
SET status = 'revoked'
WHERE id = $1;

-- name: GetSharedVaultsForUser :many
SELECT DISTINCT v.id, v.user_id, v.name, v.created_at, v.updated_at
FROM vaults v
INNER JOIN sharing_records sr ON v.id = sr.vault_id
WHERE sr.recipient_user_id = $1 AND sr.status = 'accepted'
ORDER BY v.created_at DESC;

-- name: CheckUserVaultAccess :one
SELECT 
    CASE 
        WHEN v.user_id = $2 THEN 'owner'::TEXT
        WHEN sr.recipient_user_id = $2 AND sr.status = 'accepted' THEN 'shared'::TEXT
        ELSE NULL
    END as access_level
FROM vaults v
LEFT JOIN sharing_records sr ON v.id = sr.vault_id AND sr.recipient_user_id = $2
WHERE v.id = $1
LIMIT 1;

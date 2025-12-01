-- +goose Up
-- Create devices table for storing device public keys and metadata
CREATE TABLE IF NOT EXISTS devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    device_label VARCHAR(255),
    x25519_public BYTEA NOT NULL,  -- X25519 public key for ECDH
    ed25519_public BYTEA NOT NULL, -- Ed25519 public key for signatures
    revoked_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_seen TIMESTAMP NULL
);

CREATE INDEX idx_devices_user_id ON devices(user_id);
CREATE INDEX idx_devices_revoked_at ON devices(revoked_at) WHERE revoked_at IS NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_devices_revoked_at;
DROP INDEX IF EXISTS idx_devices_user_id;
DROP TABLE IF EXISTS devices;

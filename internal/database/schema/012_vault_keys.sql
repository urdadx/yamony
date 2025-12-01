-- +goose Up
-- Create vault_keys table for storing wrapped VEKs (Vault Encryption Keys)
CREATE TABLE IF NOT EXISTS vault_keys (
    id SERIAL PRIMARY KEY,
    vault_id INTEGER NOT NULL REFERENCES vaults(id) ON DELETE CASCADE,
    wrapped_vek BYTEA NOT NULL,        -- ciphertext (wrapped by key derived from master password)
    wrap_iv BYTEA NOT NULL,            -- IV used for wrapping
    wrap_tag BYTEA NOT NULL,           -- tag for AES-GCM wrapping
    kdf_salt BYTEA NOT NULL,           -- salt used for Argon2 when creating the MK (store so client can derive MK)
    kdf_params JSONB NOT NULL,         -- stores time, memory, parallelism
    version INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(vault_id, version)
);

CREATE INDEX idx_vault_keys_vault_id ON vault_keys(vault_id);

-- +goose Down
DROP INDEX IF EXISTS idx_vault_keys_vault_id;
DROP TABLE IF EXISTS vault_keys;

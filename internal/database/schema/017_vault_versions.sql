-- +goose Up
-- Create vault_versions table for storing vault snapshot versions
CREATE TABLE IF NOT EXISTS vault_versions (
    id SERIAL PRIMARY KEY,
    vault_id INTEGER NOT NULL REFERENCES vaults(id) ON DELETE CASCADE,
    object_key TEXT NOT NULL, -- path to encrypted snapshot in object store
    mac BYTEA NULL,           -- optional authenticated mac over snapshot metadata
    created_by_device UUID NULL REFERENCES devices(id) ON DELETE SET NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_vault_versions_vault_id ON vault_versions(vault_id);
CREATE INDEX idx_vault_versions_created_at ON vault_versions(created_at);

-- +goose Down
DROP INDEX IF EXISTS idx_vault_versions_created_at;
DROP INDEX IF EXISTS idx_vault_versions_vault_id;
DROP TABLE IF EXISTS vault_versions;

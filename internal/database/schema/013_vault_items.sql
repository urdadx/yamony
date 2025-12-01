-- +goose Up
-- Create vault_items table for encrypted items with per-item encryption
CREATE TABLE IF NOT EXISTS vault_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vault_id INTEGER NOT NULL REFERENCES vaults(id) ON DELETE CASCADE,
    item_type VARCHAR(50) NOT NULL, -- login, note, card, alias, etc
    encrypted_blob BYTEA NOT NULL,  -- AES-GCM ciphertext of plaintext JSON
    iv BYTEA NOT NULL,
    tag BYTEA NOT NULL,
    meta JSONB NULL,                -- non-sensitive searchable metadata (optional: keep minimal)
    version INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_vault_items_vault_id ON vault_items(vault_id);
CREATE INDEX idx_vault_items_item_type ON vault_items(item_type);
CREATE INDEX idx_vault_items_created_at ON vault_items(created_at);

-- +goose Down
DROP INDEX IF EXISTS idx_vault_items_created_at;
DROP INDEX IF EXISTS idx_vault_items_item_type;
DROP INDEX IF EXISTS idx_vault_items_vault_id;
DROP TABLE IF EXISTS vault_items;

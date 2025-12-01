-- +goose Up
-- Create vault_attachments table for encrypted file attachments
CREATE TABLE IF NOT EXISTS vault_attachments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vault_id INTEGER NOT NULL REFERENCES vaults(id) ON DELETE CASCADE,
    item_id UUID NULL REFERENCES vault_items(id) ON DELETE SET NULL,
    object_key TEXT NOT NULL,  -- path to object in object store
    size BIGINT NOT NULL,
    iv BYTEA NOT NULL,
    tag BYTEA NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_vault_attachments_vault_id ON vault_attachments(vault_id);
CREATE INDEX idx_vault_attachments_item_id ON vault_attachments(item_id);

-- +goose Down
DROP INDEX IF EXISTS idx_vault_attachments_item_id;
DROP INDEX IF EXISTS idx_vault_attachments_vault_id;
DROP TABLE IF EXISTS vault_attachments;

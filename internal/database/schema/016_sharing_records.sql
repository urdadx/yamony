-- +goose Up
-- Create sharing_records table for encrypted wrapped keys per recipient
CREATE TABLE IF NOT EXISTS sharing_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vault_id INTEGER NOT NULL REFERENCES vaults(id) ON DELETE CASCADE,
    item_id UUID NULL REFERENCES vault_items(id) ON DELETE CASCADE,
    sender_user_id INTEGER NOT NULL REFERENCES users(id),
    recipient_user_id INTEGER NOT NULL REFERENCES users(id),
    wrapped_key BYTEA NOT NULL, -- key encrypted with ECDH-derived symmetric key
    wrap_iv BYTEA NOT NULL,
    wrap_tag BYTEA NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    accepted_at TIMESTAMP NULL
);

CREATE INDEX idx_sharing_records_vault_id ON sharing_records(vault_id);
CREATE INDEX idx_sharing_records_item_id ON sharing_records(item_id);
CREATE INDEX idx_sharing_records_recipient_user_id ON sharing_records(recipient_user_id);
CREATE INDEX idx_sharing_records_status ON sharing_records(status);

-- +goose Down
DROP INDEX IF EXISTS idx_sharing_records_status;
DROP INDEX IF EXISTS idx_sharing_records_recipient_user_id;
DROP INDEX IF EXISTS idx_sharing_records_item_id;
DROP INDEX IF EXISTS idx_sharing_records_vault_id;
DROP TABLE IF EXISTS sharing_records;

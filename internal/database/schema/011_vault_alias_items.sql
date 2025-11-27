-- +goose Up
-- Create vault_alias_items table
CREATE TABLE IF NOT EXISTS vault_alias_items (
    id SERIAL PRIMARY KEY,
    vault_id INTEGER NOT NULL REFERENCES vaults(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Core alias fields
    title VARCHAR(255) NOT NULL,
    alias_prefix VARCHAR(255) NOT NULL,
    alias_suffix VARCHAR(255) NOT NULL,
    forwards_to VARCHAR(255) NOT NULL, -- Email address
    note TEXT,
    
    -- Metadata
    is_favorite BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_used_at TIMESTAMP
);

-- Create vault_alias_attachments table for file attachments
CREATE TABLE IF NOT EXISTS vault_alias_attachments (
    id SERIAL PRIMARY KEY,
    alias_item_id INTEGER NOT NULL REFERENCES vault_alias_items(id) ON DELETE CASCADE,
    filename VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL, -- Path to encrypted file
    file_size INTEGER NOT NULL,
    mime_type VARCHAR(100),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_vault_alias_items_vault_id ON vault_alias_items(vault_id);
CREATE INDEX idx_vault_alias_items_user_id ON vault_alias_items(user_id);
CREATE INDEX idx_vault_alias_items_is_favorite ON vault_alias_items(is_favorite);
CREATE INDEX idx_vault_alias_attachments_alias_item_id ON vault_alias_attachments(alias_item_id);

-- +goose Down
DROP INDEX IF EXISTS idx_vault_alias_attachments_alias_item_id;
DROP INDEX IF EXISTS idx_vault_alias_items_is_favorite;
DROP INDEX IF EXISTS idx_vault_alias_items_user_id;
DROP INDEX IF EXISTS idx_vault_alias_items_vault_id;
DROP TABLE IF EXISTS vault_alias_attachments;
DROP TABLE IF EXISTS vault_alias_items;

-- +goose Up
-- Create vault_note_items table
CREATE TABLE IF NOT EXISTS vault_note_items (
    id SERIAL PRIMARY KEY,
    vault_id INTEGER NOT NULL REFERENCES vaults(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Core note fields
    title VARCHAR(255) NOT NULL,
    note TEXT NOT NULL,
    
    -- Metadata
    is_favorite BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_used_at TIMESTAMP
);

-- Create vault_note_attachments table for file attachments
CREATE TABLE IF NOT EXISTS vault_note_attachments (
    id SERIAL PRIMARY KEY,
    note_item_id INTEGER NOT NULL REFERENCES vault_note_items(id) ON DELETE CASCADE,
    filename VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL, -- Path to encrypted file
    file_size INTEGER NOT NULL,
    mime_type VARCHAR(100),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_vault_note_items_vault_id ON vault_note_items(vault_id);
CREATE INDEX idx_vault_note_items_user_id ON vault_note_items(user_id);
CREATE INDEX idx_vault_note_items_is_favorite ON vault_note_items(is_favorite);
CREATE INDEX idx_vault_note_attachments_note_item_id ON vault_note_attachments(note_item_id);

-- +goose Down
DROP INDEX IF EXISTS idx_vault_note_attachments_note_item_id;
DROP INDEX IF EXISTS idx_vault_note_items_is_favorite;
DROP INDEX IF EXISTS idx_vault_note_items_user_id;
DROP INDEX IF EXISTS idx_vault_note_items_vault_id;
DROP TABLE IF EXISTS vault_note_attachments;
DROP TABLE IF EXISTS vault_note_items;

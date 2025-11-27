-- +goose Up
-- Create vault_login_items table
CREATE TABLE IF NOT EXISTS vault_login_items (
    id SERIAL PRIMARY KEY,
    vault_id INTEGER NOT NULL REFERENCES vaults(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Core login fields
    title VARCHAR(255) NOT NULL,
    username VARCHAR(255),
    password_encrypted TEXT NOT NULL,
    totp_secret_encrypted TEXT, -- 2FA secret key (encrypted)
    note TEXT,
    
    -- Metadata
    is_favorite BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_used_at TIMESTAMP
);

-- Create vault_login_websites table for multiple websites per login item
CREATE TABLE IF NOT EXISTS vault_login_websites (
    id SERIAL PRIMARY KEY,
    login_item_id INTEGER NOT NULL REFERENCES vault_login_items(id) ON DELETE CASCADE,
    url VARCHAR(2048) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create vault_login_attachments table for file attachments
CREATE TABLE IF NOT EXISTS vault_login_attachments (
    id SERIAL PRIMARY KEY,
    login_item_id INTEGER NOT NULL REFERENCES vault_login_items(id) ON DELETE CASCADE,
    filename VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL, -- Path to encrypted file
    file_size INTEGER NOT NULL,
    mime_type VARCHAR(100),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_vault_login_items_vault_id ON vault_login_items(vault_id);
CREATE INDEX idx_vault_login_items_user_id ON vault_login_items(user_id);
CREATE INDEX idx_vault_login_items_is_favorite ON vault_login_items(is_favorite);
CREATE INDEX idx_vault_login_websites_login_item_id ON vault_login_websites(login_item_id);
CREATE INDEX idx_vault_login_attachments_login_item_id ON vault_login_attachments(login_item_id);

-- +goose Down
DROP INDEX IF EXISTS idx_vault_login_attachments_login_item_id;
DROP INDEX IF EXISTS idx_vault_login_websites_login_item_id;
DROP INDEX IF EXISTS idx_vault_login_items_is_favorite;
DROP INDEX IF EXISTS idx_vault_login_items_user_id;
DROP INDEX IF EXISTS idx_vault_login_items_vault_id;
DROP TABLE IF EXISTS vault_login_attachments;
DROP TABLE IF EXISTS vault_login_websites;
DROP TABLE IF EXISTS vault_login_items;

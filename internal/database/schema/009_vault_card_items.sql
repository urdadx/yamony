-- +goose Up
-- Create vault_card_items table
CREATE TABLE IF NOT EXISTS vault_card_items (
    id SERIAL PRIMARY KEY,
    vault_id INTEGER NOT NULL REFERENCES vaults(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Core card fields
    name VARCHAR(255) NOT NULL,
    card_number_encrypted TEXT NOT NULL,
    expiration_date VARCHAR(7), -- Format: MM/YYYY
    security_code_encrypted TEXT,
    pin_encrypted TEXT,
    note TEXT,
    
    -- Metadata
    is_favorite BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_used_at TIMESTAMP
);

CREATE INDEX idx_vault_card_items_vault_id ON vault_card_items(vault_id);
CREATE INDEX idx_vault_card_items_user_id ON vault_card_items(user_id);
CREATE INDEX idx_vault_card_items_is_favorite ON vault_card_items(is_favorite);

-- +goose Down
DROP INDEX IF EXISTS idx_vault_card_items_is_favorite;
DROP INDEX IF EXISTS idx_vault_card_items_user_id;
DROP INDEX IF EXISTS idx_vault_card_items_vault_id;
DROP TABLE IF EXISTS vault_card_items;

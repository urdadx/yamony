-- +goose Up
-- Create pages table
CREATE TABLE IF NOT EXISTS pages (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255),
    handle VARCHAR(255) NOT NULL UNIQUE,
    banner_image VARCHAR(500),
    image VARCHAR(500),
    bio TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_pages_user_id ON pages(user_id);
CREATE INDEX idx_pages_handle ON pages(handle);
CREATE INDEX idx_pages_is_active ON pages(is_active);

-- +goose Down
DROP INDEX IF EXISTS idx_pages_is_active;
DROP INDEX IF EXISTS idx_pages_handle;
DROP INDEX IF EXISTS idx_pages_user_id;
DROP TABLE IF EXISTS pages;

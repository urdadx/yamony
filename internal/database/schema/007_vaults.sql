-- +goose Up
-- Create vaults table
CREATE TABLE IF NOT EXISTS vaults (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    icon VARCHAR(100),
    theme VARCHAR(50),
    is_favorite BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_vaults_user_id ON vaults(user_id);
CREATE INDEX idx_vaults_is_favorite ON vaults(is_favorite);

-- +goose Down
DROP INDEX IF EXISTS idx_vaults_is_favorite;
DROP INDEX IF EXISTS idx_vaults_user_id;
DROP TABLE IF EXISTS vaults;

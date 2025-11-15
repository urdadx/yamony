-- +goose Up
-- Create preferences table
CREATE TABLE IF NOT EXISTS preferences (
    id SERIAL PRIMARY KEY,
    page_id INTEGER NOT NULL UNIQUE REFERENCES pages(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    social_icons_position VARCHAR(50) NOT NULL DEFAULT 'top',
    hide_shop BOOLEAN NOT NULL DEFAULT false,
    hide_link_branding BOOLEAN NOT NULL DEFAULT false,
    hide_share_button BOOLEAN NOT NULL DEFAULT false,
    shop_layout VARCHAR(50) NOT NULL DEFAULT 'grid',
    button_style VARCHAR(50) NOT NULL DEFAULT 'rounded-md',
    theme VARCHAR(50) DEFAULT 'Light',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_preferences_page_id ON preferences(page_id);
CREATE INDEX idx_preferences_user_id ON preferences(user_id);

-- +goose Down
DROP INDEX IF EXISTS idx_preferences_user_id;
DROP INDEX IF EXISTS idx_preferences_page_id;
DROP TABLE IF EXISTS preferences;

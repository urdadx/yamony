-- +goose Up
-- Create blocks table
CREATE TABLE IF NOT EXISTS blocks (
    id SERIAL PRIMARY KEY,
    page_id INTEGER NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
     block_order INTEGER,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- Core block properties
    title VARCHAR(500),
    description TEXT,
    layout VARCHAR(50) NOT NULL DEFAULT 'classic',
    
    -- Block type and variant
    block_type VARCHAR(100),
    variant VARCHAR(100),
    
    -- Flexible properties storage (JSON)
    properties JSONB
);

CREATE INDEX idx_blocks_page_id_block_order ON blocks(page_id, "block_order");
CREATE INDEX idx_blocks_user_id ON blocks(user_id);
CREATE INDEX idx_blocks_is_active ON blocks(is_active);
CREATE INDEX idx_blocks_block_type ON blocks(block_type);

-- +goose Down
DROP INDEX IF EXISTS idx_blocks_block_type;
DROP INDEX IF EXISTS idx_blocks_is_active;
DROP INDEX IF EXISTS idx_blocks_user_id;
DROP INDEX IF EXISTS idx_blocks_page_id_block_order;
DROP TABLE IF EXISTS blocks;

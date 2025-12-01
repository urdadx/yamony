-- +goose Up
-- Add KDF parameters to users table for master password derivation
ALTER TABLE users ADD COLUMN IF NOT EXISTS kdf_salt BYTEA;
ALTER TABLE users ADD COLUMN IF NOT EXISTS kdf_params JSONB;
ALTER TABLE users ADD COLUMN IF NOT EXISTS srp_verifier BYTEA;
ALTER TABLE users ADD COLUMN IF NOT EXISTS srp_salt BYTEA;

CREATE INDEX idx_users_email ON users(email);

-- +goose Down
DROP INDEX IF EXISTS idx_users_email;
ALTER TABLE users DROP COLUMN IF EXISTS srp_salt;
ALTER TABLE users DROP COLUMN IF EXISTS srp_verifier;
ALTER TABLE users DROP COLUMN IF EXISTS kdf_params;
ALTER TABLE users DROP COLUMN IF EXISTS kdf_salt;

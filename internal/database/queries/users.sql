-- name: CreateUser :one
INSERT INTO users (username, email, password_hash, email_verified, image)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, username, email, email_verified, image, created_at, updated_at;

-- name: GetUserByEmail :one
SELECT id, username, email, password_hash, email_verified, image, created_at, updated_at
FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT id, username, email, email_verified, image, created_at, updated_at
FROM users
WHERE id = $1;

-- name: UpdateUserEmailVerified :exec
UPDATE users
SET email_verified = $2, updated_at = NOW()
WHERE id = $1;

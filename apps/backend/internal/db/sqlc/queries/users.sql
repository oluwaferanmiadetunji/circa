-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 AND deleted_at IS NULL;

-- name: GetUserByAddress :one
SELECT * FROM users WHERE address = $1 AND deleted_at IS NULL;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 AND deleted_at IS NULL;

-- name: CreateUser :one
INSERT INTO users (full_name, email, address, display_name)
VALUES ($1, $2, $3, $4)
RETURNING *;


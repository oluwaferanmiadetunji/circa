-- name: CreateMagicLink :one
INSERT INTO magic_links (pending_signup_id, token_hash, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetMagicLinkByTokenHash :one
SELECT * FROM magic_links WHERE token_hash = $1 AND deleted_at IS NULL AND expires_at > NOW();

-- name: GetMagicLinkByPendingSignupID :one
SELECT * FROM magic_links WHERE pending_signup_id = $1 AND deleted_at IS NULL AND expires_at > NOW();

-- name: UpdateMagicLink :one
UPDATE magic_links SET expires_at = $1 WHERE id = $2 AND deleted_at IS NULL AND expires_at > NOW() RETURNING *;
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

-- name: MarkMagicLinkAsUsed :one
UPDATE magic_links
SET deleted_at = NOW()
WHERE id = $1
  AND deleted_at IS NULL
  AND expires_at > NOW()
RETURNING *;

-- name: InvalidateMagicLinksByEmail :exec
UPDATE magic_links 
SET deleted_at = NOW() 
FROM pending_signups 
WHERE magic_links.pending_signup_id = pending_signups.id 
  AND pending_signups.email = $1 
  AND magic_links.deleted_at IS NULL 
  AND pending_signups.deleted_at IS NULL 
  AND pending_signups.status = 'pending' 
  AND pending_signups.email_verified_at IS NULL;
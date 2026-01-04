-- name: CreatePendingSignup :one
INSERT INTO pending_signups (full_name, email, display_name, expires_at)
VALUES ($1, $2, $3, NOW() + INTERVAL '24 hours')
RETURNING *;

-- name: GetPendingSignupByEmail :one
SELECT * FROM pending_signups WHERE email = $1 AND deleted_at IS NULL AND status = 'pending' AND email_verified_at IS NULL AND expires_at > NOW();

-- name: GetPendingSignupByID :one
SELECT * FROM pending_signups WHERE id = $1 AND deleted_at IS NULL AND status = 'pending' AND email_verified_at IS NULL AND expires_at > NOW();

-- name: GetVerifiedPendingSignupByID :one
SELECT * FROM pending_signups WHERE id = $1 AND deleted_at IS NULL AND status = 'pending' AND email_verified_at IS NOT NULL AND expires_at > NOW();

-- name: UpdatePendingSignup :one
UPDATE pending_signups SET email_verified_at = $1 WHERE id = $2 AND deleted_at IS NULL AND status = 'pending' AND email_verified_at IS NULL AND expires_at > NOW() RETURNING *;

-- name: InvalidatePendingSignupsByEmail :exec
UPDATE pending_signups SET deleted_at = NOW() WHERE email = $1 AND deleted_at IS NULL AND status = 'pending' AND email_verified_at IS NULL;

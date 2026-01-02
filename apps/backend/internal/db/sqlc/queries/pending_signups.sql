-- name: CreatePendingSignup :one
INSERT INTO pending_signups (full_name, email, display_name)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetPendingSignupByEmail :one
SELECT * FROM pending_signups WHERE email = $1 AND deleted_at IS NULL AND status = 'pending' AND email_verified_at IS NULL AND expires_at > NOW();

-- name: GetPendingSignupByID :one
SELECT * FROM pending_signups WHERE id = $1 AND deleted_at IS NULL AND status = 'pending' AND email_verified_at IS NULL AND expires_at > NOW();

-- name: UpdatePendingSignup :one
UPDATE pending_signups SET email_verified_at = $1 WHERE id = $2 AND deleted_at IS NULL AND status = 'pending' AND email_verified_at IS NULL AND expires_at > NOW() RETURNING *;
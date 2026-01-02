-- name: CreateJob :one
INSERT INTO jobs (type, payload, max_retries, scheduled_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetNextPendingJob :one
SELECT * FROM jobs 
WHERE status = 'pending' 
  AND deleted_at IS NULL 
  AND scheduled_at <= NOW()
ORDER BY scheduled_at ASC, created_at ASC
LIMIT 1
FOR UPDATE SKIP LOCKED;

-- name: UpdateJobStatus :one
UPDATE jobs 
SET status = $1, 
    processed_at = CASE WHEN $1 = 'completed' OR $1 = 'failed' THEN NOW() ELSE processed_at END,
    error_message = $2,
    updated_at = NOW()
WHERE id = $3 AND deleted_at IS NULL
RETURNING *;

-- name: IncrementJobRetry :one
UPDATE jobs 
SET retry_count = retry_count + 1,
    status = CASE 
        WHEN retry_count + 1 >= max_retries THEN 'failed'
        ELSE 'pending'
    END,
    error_message = $1,
    scheduled_at = NOW() + INTERVAL '5 minutes' * (retry_count + 1),
    updated_at = NOW()
WHERE id = $2 AND deleted_at IS NULL
RETURNING *;

-- name: GetJobByID :one
SELECT * FROM jobs WHERE id = $1 AND deleted_at IS NULL;


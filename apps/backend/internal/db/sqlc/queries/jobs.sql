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
UPDATE jobs j
SET 
    status = v.status_val,
    processed_at = CASE 
        WHEN v.status_val IN ('completed', 'failed')
        THEN NOW()
        ELSE j.processed_at
    END,
    error_message = v.error_msg,
    updated_at = NOW()
FROM (VALUES (sqlc.arg(status)::varchar, sqlc.arg(error_message), sqlc.arg(id)::uuid)) AS v(status_val, error_msg, job_id)
WHERE j.id = v.job_id AND j.deleted_at IS NULL
RETURNING j.*;

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


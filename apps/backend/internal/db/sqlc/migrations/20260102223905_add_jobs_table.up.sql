CREATE TABLE
    jobs (
        "id" UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        "type" VARCHAR NOT NULL,
        "payload" JSONB NOT NULL,
        "status" VARCHAR NOT NULL DEFAULT 'pending',
        "retry_count" INTEGER NOT NULL DEFAULT 0,
        "max_retries" INTEGER NOT NULL DEFAULT 3,
        "error_message" TEXT,
        "scheduled_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        "processed_at" TIMESTAMPTZ,
        "created_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        "updated_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        "deleted_at" TIMESTAMP
    );

CREATE INDEX idx_jobs_status_scheduled ON jobs (status, scheduled_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_jobs_type ON jobs (type) WHERE deleted_at IS NULL AND status = 'pending';


CREATE TABLE
    pending_signups (
        "id" UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        "full_name" VARCHAR,
        "email" VARCHAR,
        "display_name" TEXT,
        "email_verified_at" TIMESTAMP,
        "status" TEXT NOT NULL DEFAULT 'pending',
        "expires_at" TIMESTAMP,
        "created_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        "updated_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        "deleted_at" TIMESTAMP
    );
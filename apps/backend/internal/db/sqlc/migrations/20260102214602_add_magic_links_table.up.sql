CREATE TABLE
    magic_links (
        "id" UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        "pending_signup_id" UUID NOT NULL REFERENCES pending_signups (id),
        "token_hash" VARCHAR NOT NULL,
        "expires_at" TIMESTAMP,
        "created_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        "updated_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        "deleted_at" TIMESTAMP
    );
CREATE TABLE
  users (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    "full_name" VARCHAR,
    "email" VARCHAR,
    "address" VARCHAR UNIQUE NOT NULL,
    "display_name" TEXT,
    "avatar_url" TEXT,
    "created_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP
  );
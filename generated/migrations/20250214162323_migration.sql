-- +goose Up
-- Migration script generated from GraphQL schema (incremental)

CREATE TABLE "user_profile" (
    "id" TEXT PRIMARY KEY,
    "bio" TEXT,
    "avatar_url" TEXT,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE "user" (
    "id" TEXT PRIMARY KEY,
    "name" TEXT NOT NULL,
    "email" TEXT,
    "created_date" TIMESTAMPTZ,
    "is_active" BOOLEAN,
    "profile_id" TEXT,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY ("profile_id") REFERENCES "user_profile"("id") ON DELETE CASCADE
);

CREATE INDEX "idx_user_name" ON "user"("name");
CREATE INDEX "idx_user_email" ON "user"("email");

CREATE TABLE "post" (
    "id" TEXT PRIMARY KEY,
    "title" TEXT NOT NULL,
    "content" TEXT,
    "published_date" TIMESTAMPTZ,
    "author_id" TEXT,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY ("author_id") REFERENCES "user"("id") ON DELETE CASCADE
);

CREATE INDEX "idx_post_author" ON "post"("author_id");

CREATE TABLE "collection" (
    "id" TEXT PRIMARY KEY,
    "address" TEXT NOT NULL,
    "type" TEXT,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX "idx_collection_address" ON "collection"("address");

CREATE TABLE "transfer" (
    "id" TEXT PRIMARY KEY,
    "from" TEXT NOT NULL,
    "to" TEXT NOT NULL,
    "amount" NUMERIC,
    "timestamp" TIMESTAMPTZ,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE "post" ADD COLUMN "user_id" TEXT;
ALTER TABLE "post" ADD CONSTRAINT "fk_post_user" 
					FOREIGN KEY ("user_id") REFERENCES "user"("id") ON DELETE CASCADE;
CREATE INDEX "idx_post_user_id" ON "post"("user_id");


-- +goose Down

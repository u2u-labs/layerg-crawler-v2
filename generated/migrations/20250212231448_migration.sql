-- +goose Up
-- Migration script generated from GraphQL schema (incremental)

CREATE TABLE "user_profile" (
    "id" SERIAL PRIMARY KEY,
    "bio" TEXT,
    "avatar_url" TEXT
);


CREATE TABLE "user" (
    "id" SERIAL PRIMARY KEY,
    "name" TEXT NOT NULL,
    "email" TEXT,
    "created_date" TIMESTAMPTZ,
    "is_active" BOOLEAN,
    "profile_id" INTEGER,
    FOREIGN KEY ("profile_id") REFERENCES "user_profile"("id")
);

CREATE INDEX "idx_user_name" ON "user"("name");
CREATE INDEX "idx_user_email" ON "user"("email");

CREATE TABLE "post" (
    "id" SERIAL PRIMARY KEY,
    "title" TEXT NOT NULL,
    "content" TEXT,
    "published_date" TIMESTAMPTZ,
    "author_id" INTEGER,
    FOREIGN KEY ("author_id") REFERENCES "user"("id")
);

CREATE INDEX "idx_post_author" ON "post"("author_id");

CREATE TABLE "collection" (
    "id" SERIAL PRIMARY KEY,
    "address" TEXT NOT NULL,
    "type" TEXT
);

CREATE INDEX "idx_collection_address" ON "collection"("address");

CREATE TABLE "transfer" (
    "id" SERIAL PRIMARY KEY,
    "from" TEXT NOT NULL,
    "to" TEXT NOT NULL,
    "amount" DECIMAL,
    "timestamp" TIMESTAMPTZ
);



-- +goose Down

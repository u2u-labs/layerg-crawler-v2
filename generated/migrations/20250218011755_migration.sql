-- +goose Up
-- Migration script generated from GraphQL schema (incremental)

CREATE TABLE "user" (
    "id" TEXT PRIMARY KEY,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE "item" (
    "id" TEXT PRIMARY KEY,
    "token_id" NUMERIC NOT NULL,
    "token_uri" TEXT NOT NULL,
    "owner_id" TEXT NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY ("owner_id") REFERENCES "user"("id") ON DELETE CASCADE
);


CREATE TABLE "metadata_update_record" (
    "id" TEXT PRIMARY KEY,
    "token_id" NUMERIC NOT NULL,
    "actor_id" TEXT NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY ("actor_id") REFERENCES "user"("id") ON DELETE CASCADE
);


ALTER TABLE "item" ADD COLUMN "user_id" TEXT;
ALTER TABLE "item" ADD CONSTRAINT "fk_item_user" 
					FOREIGN KEY ("user_id") REFERENCES "user"("id") ON DELETE CASCADE;
CREATE INDEX "idx_item_user_id" ON "item"("user_id");


-- +goose Down

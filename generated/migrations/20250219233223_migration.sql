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
    "standard" TEXT NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE "balance" (
    "id" TEXT PRIMARY KEY,
    "item_id" TEXT NOT NULL,
    "owner_id" TEXT NOT NULL,
    "value" NUMERIC NOT NULL,
    "updated_at" NUMERIC NOT NULL,
    "contract" TEXT NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY ("item_id") REFERENCES "item"("id") ON DELETE CASCADE,
    FOREIGN KEY ("owner_id") REFERENCES "user"("id") ON DELETE CASCADE
);


CREATE TABLE "metadata_update_record" (
    "id" TEXT PRIMARY KEY,
    "token_id" NUMERIC NOT NULL,
    "actor_id" TEXT NOT NULL,
    "timestamp" NUMERIC NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY ("actor_id") REFERENCES "user"("id") ON DELETE CASCADE
);



-- +goose Down

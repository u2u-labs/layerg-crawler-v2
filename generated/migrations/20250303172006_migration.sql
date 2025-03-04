-- +goose Up
-- Migration script generated from GraphQL schema (incremental)

CREATE TABLE "value" (
    "id" TEXT PRIMARY KEY,
    "value" NUMERIC NOT NULL,
    "sender" TEXT NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);



-- +goose Down
DROP TABLE IF EXISTS "value" CASCADE;

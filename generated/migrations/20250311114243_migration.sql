-- +goose Up
-- Migration script generated from GraphQL schema (incremental)

CREATE TYPE "ItemStandard" AS ENUM (
    'ERC721',
    'ERC1155'
);

CREATE TABLE "item" (
    "id" TEXT PRIMARY KEY,
    "token_id" NUMERIC NOT NULL,
    "token_uri" TEXT NOT NULL,
    "standard" TEXT NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE "user" (
    "id" TEXT PRIMARY KEY,
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

CREATE UNIQUE INDEX "idx_balance_item" ON "balance"("item_id");
CREATE INDEX "idx_balance_value" ON "balance"("value");
CREATE INDEX "idx_composite_balance_0" ON "balance"("item_id", "owner_id");
CREATE INDEX "idx_composite_balance_1" ON "balance"("item_id", "value");

CREATE TABLE "metadata_update_record" (
    "id" TEXT PRIMARY KEY,
    "token_id" NUMERIC NOT NULL,
    "actor_id" TEXT NOT NULL,
    "timestamp" NUMERIC NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY ("actor_id") REFERENCES "user"("id") ON DELETE CASCADE
);



-- +goose Down
DROP TABLE IF EXISTS "metadata_update_record" CASCADE;
DROP TABLE IF EXISTS "balance" CASCADE;
DROP TABLE IF EXISTS "user" CASCADE;
DROP TABLE IF EXISTS "item" CASCADE;
DROP TYPE IF EXISTS "ItemStandard" CASCADE;

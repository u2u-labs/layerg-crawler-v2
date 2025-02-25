-- name: CreateItem :one
INSERT INTO "item" ("id", "token_id", "token_uri", "owner_id", "contract") VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetItem :one
SELECT * FROM "item" WHERE id = $1;

-- name: ListItem :many
SELECT * FROM "item";

-- name: UpdateItem :one
UPDATE "item" SET "token_id" = $2, "token_uri" = $3, "owner_id" = $4, "contract" = $5 WHERE id = $1 RETURNING *;

-- name: DeleteItem :exec
DELETE FROM "item" WHERE id = $1;

-- name: CreateMetadataUpdateRecord :one
INSERT INTO "metadata_update_record" ("id", "item_id", "actor_id") VALUES ($1, $2, $3) RETURNING *;

-- name: GetMetadataUpdateRecord :one
SELECT * FROM "metadata_update_record" WHERE id = $1;

-- name: ListMetadataUpdateRecord :many
SELECT * FROM "metadata_update_record";

-- name: UpdateMetadataUpdateRecord :one
UPDATE "metadata_update_record" SET "item_id" = $2, "actor_id" = $3 WHERE id = $1 RETURNING *;

-- name: DeleteMetadataUpdateRecord :exec
DELETE FROM "metadata_update_record" WHERE id = $1;

-- name: CreateUser :one
INSERT INTO "user" ("id") VALUES ($1) RETURNING *;

-- name: GetUser :one
SELECT * FROM "user" WHERE id = $1;

-- name: ListUser :many
SELECT * FROM "user";

-- name: UpdateUser :exec
-- Skip update query generation as there are no updateable fields

-- name: DeleteUser :exec
DELETE FROM "user" WHERE id = $1;


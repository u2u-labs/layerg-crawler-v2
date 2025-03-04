-- name: CreateValue :one
INSERT INTO "value" ("id", "value", "sender") VALUES ($1, $2, $3) RETURNING *;

-- name: GetValue :one
SELECT * FROM "value" WHERE id = $1;

-- name: ListValue :many
SELECT * FROM "value";

-- name: UpdateValue :one
UPDATE "value" SET "value" = $2, "sender" = $3 WHERE id = $1 RETURNING *;

-- name: DeleteValue :exec
DELETE FROM "value" WHERE id = $1;


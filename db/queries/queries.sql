-- name: CreateUser :one
INSERT INTO "user" ("name", "email", "created_date", "is_active", "profile_id") VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetUser :one
SELECT * FROM "user" WHERE id = $1;

-- name: ListUser :many
SELECT * FROM "user";

-- name: UpdateUser :one
UPDATE "user" SET "name" = $2, "email" = $3, "created_date" = $4, "is_active" = $5, "profile_id" = $6 WHERE id = $1 RETURNING *;

-- name: DeleteUser :exec
DELETE FROM "user" WHERE id = $1;

-- name: CreateUserProfile :one
INSERT INTO "user_profile" ("bio", "avatar_url") VALUES ($1, $2) RETURNING *;

-- name: GetUserProfile :one
SELECT * FROM "user_profile" WHERE id = $1;

-- name: ListUserProfile :many
SELECT * FROM "user_profile";

-- name: UpdateUserProfile :one
UPDATE "user_profile" SET "bio" = $2, "avatar_url" = $3 WHERE id = $1 RETURNING *;

-- name: DeleteUserProfile :exec
DELETE FROM "user_profile" WHERE id = $1;

-- name: CreatePost :one
INSERT INTO "post" ("title", "content", "published_date", "author_id") VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetPost :one
SELECT * FROM "post" WHERE id = $1;

-- name: ListPost :many
SELECT * FROM "post";

-- name: UpdatePost :one
UPDATE "post" SET "title" = $2, "content" = $3, "published_date" = $4, "author_id" = $5 WHERE id = $1 RETURNING *;

-- name: DeletePost :exec
DELETE FROM "post" WHERE id = $1;

-- name: CreateCollection :one
INSERT INTO "collection" ("address", "type") VALUES ($1, $2) RETURNING *;

-- name: GetCollection :one
SELECT * FROM "collection" WHERE id = $1;

-- name: ListCollection :many
SELECT * FROM "collection";

-- name: UpdateCollection :one
UPDATE "collection" SET "address" = $2, "type" = $3 WHERE id = $1 RETURNING *;

-- name: DeleteCollection :exec
DELETE FROM "collection" WHERE id = $1;

-- name: CreateTransfer :one
INSERT INTO "transfer" ("from", "to", "amount", "timestamp") VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetTransfer :one
SELECT * FROM "transfer" WHERE id = $1;

-- name: ListTransfer :many
SELECT * FROM "transfer";

-- name: UpdateTransfer :one
UPDATE "transfer" SET "from" = $2, "to" = $3, "amount" = $4, "timestamp" = $5 WHERE id = $1 RETURNING *;

-- name: DeleteTransfer :exec
DELETE FROM "transfer" WHERE id = $1;


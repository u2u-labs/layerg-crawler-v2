-- Add a new query to get or create user
-- name: GetOrCreateUser :one
INSERT INTO "user" ("id") 
VALUES ($1) 
ON CONFLICT (id) DO UPDATE SET id = EXCLUDED.id
RETURNING *;


-- name: GetItemByTokenIdAndContract :one
SELECT * FROM "item" WHERE token_id = $1 AND contract = $2;

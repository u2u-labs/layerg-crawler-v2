-- Add a new query to get or create user
-- name: GetOrCreateUser :one
INSERT INTO "user" ("id") 
VALUES ($1) 
ON CONFLICT (id) DO UPDATE SET id = EXCLUDED.id
RETURNING *;

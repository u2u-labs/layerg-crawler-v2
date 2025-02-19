-- Add a new query to get or create user
-- name: GetOrCreateUser :one
INSERT INTO "user" ("id") 
VALUES ($1) 
ON CONFLICT (id) DO UPDATE SET id = EXCLUDED.id
RETURNING *;

-- name: GetUserBalance :one
SELECT id, value FROM balance 
WHERE owner_id = $1 AND item_id = $2
LIMIT 1;

-- name: UpsertBalance :one
INSERT INTO balance (
    id,
    item_id,
    owner_id,
    value,
    updated_at,
    contract
)
VALUES (
    $1, -- id (UUID)
    $2, -- item_id
    $3, -- owner_id
    $4, -- value
    $5, -- updated_at (block timestamp)
    $6  -- contract address
)
ON CONFLICT (id) 
DO UPDATE SET 
    value = EXCLUDED.value,
    updated_at = EXCLUDED.updated_at
RETURNING *;

-- name: GetItemByTokenId :one
SELECT id, token_id, token_uri, standard, created_at FROM "item" WHERE token_id = $1;
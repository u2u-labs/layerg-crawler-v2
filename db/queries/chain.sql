-- name: GetChainById :one
SELECT * FROM chains WHERE id = $1;

-- name: GetAllChain :many
SELECT * FROM chains;

-- name: UpdateChainLatestBlock :exec
UPDATE chains
SET
    latest_block = $2
WHERE
    id = $1;
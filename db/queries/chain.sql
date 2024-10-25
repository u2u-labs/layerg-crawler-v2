-- name: GetChainById :one
SELECT * FROM chains WHERE id = $1;

-- name: GetAllChain :many
SELECT * FROM chains;




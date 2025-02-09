-- name: GetAllChain :many
SELECT * FROM chains;

-- name: GetChainById :one
SELECT * FROM chains WHERE id = $1;

-- name: GetPaginatedAssetsByChainId :many
SELECT * FROM assets 
WHERE chain_id = $1 
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;


-- name: UpdateChainLatestBlock :exec
UPDATE chains 
SET latest_block = $2 
WHERE id = $1;

-- name: CreateAsset :one
INSERT INTO assets (
    id, chain_id, contract_address, initial_block
) VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetAssetByAddress :one
SELECT * FROM assets 
WHERE chain_id = $1 AND contract_address = $2;

-- name: CreateOnchainHistory :one
INSERT INTO onchain_histories (
    "from", "to", chain_id, asset_id, tx_hash, receipt, event_type, timestamp
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetOnchainHistoriesByAsset :many
SELECT * FROM onchain_histories 
WHERE asset_id = $1 
ORDER BY timestamp DESC 
LIMIT $2;

-- name: CreateChain :one
INSERT INTO chains (
    id, chain, name, rpc_url, chain_id, explorer, latest_block, block_time
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *; 
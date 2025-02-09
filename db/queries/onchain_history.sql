-- name: CreateOnchainHistory :one
INSERT INTO onchain_histories (
    "from", 
    "to", 
    chain_id, 
    asset_id, 
    tx_hash, 
    receipt, 
    event_type, 
    timestamp,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, 
    CURRENT_TIMESTAMP, 
    CURRENT_TIMESTAMP
) RETURNING *;

-- name: GetOnchainHistoryByTxHash :one
SELECT * FROM onchain_histories 
WHERE tx_hash = $1;

-- name: GetOnchainHistoriesByAddress :many
SELECT * FROM onchain_histories 
WHERE "from" = $1 OR "to" = $1
ORDER BY timestamp DESC
LIMIT $2;

-- name: GetOnchainHistoriesByAssetAndAddress :many
SELECT * FROM onchain_histories 
WHERE asset_id = $1 AND ("from" = $2 OR "to" = $2)
ORDER BY timestamp DESC
LIMIT $3; 
-- name: GetAssetByChainIdAndContractAddress :one
SELECT * FROM assets 
WHERE chain_id = $1 
AND collection_address = $2;

-- name: GetAssetByChainId :many
SELECT * FROM assets 
WHERE chain_id = $1;

-- name: AddNewAsset :exec
INSERT INTO assets (
    id, chain_id, collection_address, type, decimal_data, initial_block, last_updated
)
VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;
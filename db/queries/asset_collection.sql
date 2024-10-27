-- name: GetAssetByChainIdAndContractAddress :one
SELECT * FROM assets 
WHERE chain_id = $1 
AND collection_address = $2;

-- name: GetAssetByChainId :many
SELECT * FROM assets
WHERE chain_id = $1;
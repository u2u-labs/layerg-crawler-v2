-- name: GetPaginated1155AssetByAssetId :many
SELECT * FROM erc_1155_collection_assets 
WHERE asset_id = $1
LIMIT $2 OFFSET $3;

-- name: Count1155AssetByAssetId :one
SELECT COUNT(*) FROM erc_1155_collection_assets 
WHERE asset_id = $1;

-- name: Get1155AssetByAssetIdAndTokenId :one
SELECT * FROM erc_1155_collection_assets
WHERE
    asset_id = $1
    AND token_id = $2;

-- name: GetPaginated1155AssetByOwnerAddress :many
SELECT * FROM erc_1155_collection_assets
WHERE
    owner = $1
LIMIT $2 OFFSET $3;

-- name: Count1155AssetByOwner :one
SELECT COUNT(*) FROM erc_1155_collection_assets 
WHERE owner = $1;

-- name: Add1155Asset :exec
INSERT INTO
    erc_1155_collection_assets (asset_id, chain_id, token_id, owner, balance, attributes)
VALUES (
    $1, $2, $3, $4, $5, $6
) ON CONFLICT ON CONSTRAINT UC_ERC1155 DO UPDATE SET
    balance = $5,
    attributes = $6
    
RETURNING *;

-- name: Update1155Asset :exec
UPDATE erc_1155_collection_assets
SET
    owner = $2 
WHERE 
    id = $1;

-- name: Delete1155Asset :exec
DELETE 
FROM erc_1155_collection_assets
WHERE
    id = $1;

-- name: Update1155AssetTotalSupply :exec
WITH total_balance AS (
  SELECT SUM(balance) AS total_supply
  FROM erc_1155_collection_assets
  WHERE asset_id = $1
  AND token_id = $2
)
UPDATE erc_1155_collection_assets erc1155
SET total_supply = (SELECT total_supply FROM total_balance)
WHERE erc1155.asset_id = $1
AND erc1155.token_id = $2;
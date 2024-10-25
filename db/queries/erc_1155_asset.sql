-- name: Get1155AssetByAssetId :many
SELECT * FROM erc_1155_collection_assets WHERE asset_id = $1;

-- name: Get1155AssetByAssetIdAndTokenId :one
SELECT * FROM erc_1155_collection_assets
WHERE
    asset_id = $1
    AND token_id = $2;

-- name: Get1155AssetByOwner :many
SELECT * FROM erc_1155_collection_assets
WHERE
    owner = $1;

-- name: Add1155Asset :exec
INSERT INTO
    erc_1155_collection_assets (asset_id, token_id, owner, balance, attributes)
VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

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
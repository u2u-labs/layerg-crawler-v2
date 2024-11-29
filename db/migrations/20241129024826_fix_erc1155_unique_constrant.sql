-- +goose Up
-- +goose StatementBegin

DROP INDEX IF EXISTS UC_ERC1155 CASCADE;

ALTER TABLE erc_1155_collection_assets 
ADD COLUMN total_supply DECIMAL(78, 0) NOT NULL DEFAULT 0;

ALTER TABLE erc_1155_collection_assets 
ADD CONSTRAINT UC_ERC1155 UNIQUE (asset_id, chain_id, token_id, owner);


-- Update total supply for all assets
WITH aggregated_totals AS (
  SELECT 
    asset_id,
    token_id,
    SUM(balance) AS total_supply
  FROM 
    erc_1155_collection_assets
  GROUP BY 
    asset_id, token_id
)
UPDATE erc_1155_collection_assets AS e
SET total_supply = a.total_supply
FROM aggregated_totals AS a
WHERE 
    e.asset_id = a.asset_id 
    AND e.token_id = a.token_id;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd

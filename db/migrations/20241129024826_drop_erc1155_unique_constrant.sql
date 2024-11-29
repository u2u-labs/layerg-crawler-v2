-- +goose Up
-- +goose StatementBegin

DROP INDEX IF EXISTS UC_ERC1155 CASCADE;

ALTER TABLE erc_1155_collection_assets 
ADD COLUMN total_supply DECIMAL(78, 0) NOT NULL DEFAULT 0;

ALTER TABLE erc_1155_collection_assets 
ADD CONSTRAINT UC_ERC1155 UNIQUE (asset_id, chain_id, token_id, owner);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd

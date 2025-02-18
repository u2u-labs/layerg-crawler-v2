-- +goose Up
-- +goose StatementBegin

-- Create table for chain information.    
CREATE TABLE chains (
    id INT PRIMARY KEY,
    chain VARCHAR NOT NULL,
    name VARCHAR NOT NULL,
    rpc_url VARCHAR NOT NULL,
    chain_id BIGINT NOT NULL,
    explorer VARCHAR NOT NULL,
    latest_block BIGINT NOT NULL,
    block_time INT NOT NULL
);

-- Create table for assets, linking collection addresses to a chain.
CREATE TABLE assets (
    id VARCHAR PRIMARY KEY,
    chain_id INT NOT NULL,
    contract_address VARCHAR(42) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    initial_block BIGINT,
    last_updated TIMESTAMP,
    FOREIGN KEY (chain_id) REFERENCES chains (id),
    CONSTRAINT UC_ASSET_COLLECTION UNIQUE (chain_id, contract_address)
);

-- Create table for onchain histories.
-- Stores a JSONB representation of the receipt along with the tx hash and event info.
CREATE TABLE onchain_histories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "from" VARCHAR(42) NOT NULL,
    "to" VARCHAR(42) NOT NULL,
    chain_id INT NOT NULL,
    asset_id VARCHAR NOT NULL,
    tx_hash VARCHAR(66) NOT NULL,
    receipt JSONB NOT NULL,
    event_type TEXT,
    timestamp TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_onchain_history_chain FOREIGN KEY (chain_id) REFERENCES chains(id),
    CONSTRAINT fk_onchain_history_asset FOREIGN KEY (asset_id) REFERENCES assets(id)
);

-- INSERT INTO chains (id, chain, name, rpc_url, chain_id, explorer, latest_block, block_time)
-- VALUES (1, 'U2U', 'Nebulas Testnet', 'https://rpc-nebulas-testnet.uniultra.xyz', 2484, 'https://testnet.u2uscan.xyz/', 47984307, 500);
-- Create enumeration type for backfill crawler status.
-- CREATE TYPE crawler_status AS ENUM ('CRAWLING', 'CRAWLED');


-- Create table for backfill crawlers.
-- Associates a backfill worker with a specific asset (via asset_id) and stores its current block and status.
-- CREATE TABLE backfill_crawlers (
--     chain_id INT NOT NULL,
--     asset_id INT NOT NULL,
--     current_block BIGINT NOT NULL,
--     status crawler_status DEFAULT 'CRAWLING' NOT NULL,
--     created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
--     PRIMARY KEY (chain_id, asset_id),
--     CONSTRAINT fk_bf_chain FOREIGN KEY (chain_id) REFERENCES chains(id),
--     CONSTRAINT fk_bf_asset FOREIGN KEY (asset_id) REFERENCES assets(id)
-- );

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- DROP TABLE IF EXISTS backfill_crawlers;
DROP TYPE IF EXISTS crawler_status;
DROP TABLE IF EXISTS onchain_histories;
DROP TABLE IF EXISTS assets;
DROP TABLE IF EXISTS chains;

-- +goose StatementEnd 
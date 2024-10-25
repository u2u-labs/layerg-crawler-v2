-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    tracking_contracts (
        chain_id INT NOT NULL,
        contract_address VARCHAR NOT NULL,
        PRIMARY KEY (chain_id, contract_address)
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tracking_contracts;

-- +goose StatementEnd
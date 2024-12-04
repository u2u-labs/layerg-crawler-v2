-- +goose Up
-- +goose StatementBegin
CREATE TYPE crawler_status AS ENUM ('CRAWLING', 'CRAWLED');

CREATE TABLE backfill_crawlers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    chain_id INT NOT NULL,
    collection_address VARCHAR NOT NULL,
    current_block BIGINT NOT NULL,
    status crawler_status DEFAULT 'CRAWLING',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS backfill_crawlers;

DROP TYPE IF EXISTS crawler_status;
-- +goose StatementEnd

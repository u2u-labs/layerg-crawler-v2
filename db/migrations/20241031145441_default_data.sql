-- +goose Up
-- +goose StatementBegin
INSERT INTO chains (id, chain, name, rpc_url, chain_id, explorer, latest_block, block_time)
VALUES (1, 'U2U', 'Nebulas Testnet', 'https://rpc-nebulas-testnet.uniultra.xyz', 2484, 'https://testnet.u2uscan.xyz/', 40573897, 500);

-- INSERT INTO chains (id, chain, name, rpc_url, chain_id, explorer, latest_block, block_time)
-- VALUES (2, 'U2U', 'Solaris Mainnet', 'https://rpc-mainnet.uniultra.xyz', 39, 'https://u2uscan.xyz/', 20233972, 2000);

INSERT INTO assets (id, chain_id, collection_address, type, created_at, updated_at, decimal_data, initial_block, last_updated)
VALUES ('1:0xdFAe88F8610a038AFcDF47A5BC77C0963C65087c', 1, '0xdFAe88F8610a038AFcDF47A5BC77C0963C65087c', 'ERC20', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 18, 0, CURRENT_TIMESTAMP);

INSERT INTO assets (id, chain_id, collection_address, type, created_at, updated_at, decimal_data, initial_block, last_updated)
VALUES ('2:0xC5f15624b4256C1206e4BB93f2CCc9163A75b703', 1, '0xC5f15624b4256C1206e4BB93f2CCc9163A75b703', 'ERC20', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 18, 0, CURRENT_TIMESTAMP);

INSERT INTO assets (id, chain_id, collection_address, type, created_at, updated_at, decimal_data, initial_block, last_updated)
VALUES ('1:0x0091BD12166d29539Db6bb37FB79670779aBf266', 1, '0x0091BD12166d29539Db6bb37FB79670779aBf266', 'ERC721', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 0, 0, CURRENT_TIMESTAMP);

INSERT INTO assets (id, chain_id, collection_address, type, created_at, updated_at, decimal_data, initial_block, last_updated)
VALUES ('1:0x9E87754dAB31dAD057DCDF233000F71fF55fA37f', 1, '0x9E87754dAB31dAD057DCDF233000F71fF55fA37f', 'ERC1155', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 0, 0, CURRENT_TIMESTAMP);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd

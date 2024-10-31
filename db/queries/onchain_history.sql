-- name: AddOnchainTransaction :exec
INSERT INTO 
    onchain_histories("from","to",asset_id,token_id,amount,tx_hash,timestamp)
VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

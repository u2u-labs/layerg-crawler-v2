-- name: GetAllTrackingContractOnChain :many
SELECT * FROM tracking_contracts WHERE chain_id = $1;




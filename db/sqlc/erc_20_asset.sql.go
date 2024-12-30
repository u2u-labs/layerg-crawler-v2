// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: erc_20_asset.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const add20Asset = `-- name: Add20Asset :exec
INSERT INTO
    erc_20_collection_assets (asset_id, chain_id, owner, balance)
VALUES (
    $1, $2, $3, $4
) ON CONFLICT (owner) DO UPDATE SET
    balance = $4
RETURNING id, chain_id, asset_id, owner, balance, created_at, updated_at
`

type Add20AssetParams struct {
	AssetID string `json:"assetId"`
	ChainID int32  `json:"chainId"`
	Owner   string `json:"owner"`
	Balance string `json:"balance"`
}

func (q *Queries) Add20Asset(ctx context.Context, arg Add20AssetParams) error {
	_, err := q.db.ExecContext(ctx, add20Asset,
		arg.AssetID,
		arg.ChainID,
		arg.Owner,
		arg.Balance,
	)
	return err
}

const count20AssetByAssetId = `-- name: Count20AssetByAssetId :one
SELECT COUNT(*) FROM erc_20_collection_assets 
WHERE asset_id = $1
`

func (q *Queries) Count20AssetByAssetId(ctx context.Context, assetID string) (int64, error) {
	row := q.db.QueryRowContext(ctx, count20AssetByAssetId, assetID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const count20AssetByOwner = `-- name: Count20AssetByOwner :one
SELECT COUNT(*) FROM erc_20_collection_assets 
WHERE owner = $1
`

func (q *Queries) Count20AssetByOwner(ctx context.Context, owner string) (int64, error) {
	row := q.db.QueryRowContext(ctx, count20AssetByOwner, owner)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const count20AssetHoldersByAssetId = `-- name: Count20AssetHoldersByAssetId :one
SELECT COUNT(DISTINCT(owner)) FROM erc_20_collection_assets 
WHERE asset_id = $1
`

func (q *Queries) Count20AssetHoldersByAssetId(ctx context.Context, assetID string) (int64, error) {
	row := q.db.QueryRowContext(ctx, count20AssetHoldersByAssetId, assetID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const delete20Asset = `-- name: Delete20Asset :exec
DELETE 
FROM erc_20_collection_assets
WHERE
    id = $1
`

func (q *Queries) Delete20Asset(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, delete20Asset, id)
	return err
}

const get20AssetByAssetIdAndTokenId = `-- name: Get20AssetByAssetIdAndTokenId :one
SELECT id, chain_id, asset_id, owner, balance, created_at, updated_at FROM erc_20_collection_assets
WHERE
    asset_id = $1
    AND owner = $2
`

type Get20AssetByAssetIdAndTokenIdParams struct {
	AssetID string `json:"assetId"`
	Owner   string `json:"owner"`
}

func (q *Queries) Get20AssetByAssetIdAndTokenId(ctx context.Context, arg Get20AssetByAssetIdAndTokenIdParams) (Erc20CollectionAsset, error) {
	row := q.db.QueryRowContext(ctx, get20AssetByAssetIdAndTokenId, arg.AssetID, arg.Owner)
	var i Erc20CollectionAsset
	err := row.Scan(
		&i.ID,
		&i.ChainID,
		&i.AssetID,
		&i.Owner,
		&i.Balance,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getPaginated20AssetByAssetId = `-- name: GetPaginated20AssetByAssetId :many
SELECT id, chain_id, asset_id, owner, balance, created_at, updated_at FROM erc_20_collection_assets 
WHERE asset_id = $1
LIMIT $2 OFFSET $3
`

type GetPaginated20AssetByAssetIdParams struct {
	AssetID string `json:"assetId"`
	Limit   int32  `json:"limit"`
	Offset  int32  `json:"offset"`
}

func (q *Queries) GetPaginated20AssetByAssetId(ctx context.Context, arg GetPaginated20AssetByAssetIdParams) ([]Erc20CollectionAsset, error) {
	rows, err := q.db.QueryContext(ctx, getPaginated20AssetByAssetId, arg.AssetID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Erc20CollectionAsset
	for rows.Next() {
		var i Erc20CollectionAsset
		if err := rows.Scan(
			&i.ID,
			&i.ChainID,
			&i.AssetID,
			&i.Owner,
			&i.Balance,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPaginated20AssetByOwnerAddress = `-- name: GetPaginated20AssetByOwnerAddress :many
SELECT id, chain_id, asset_id, owner, balance, created_at, updated_at FROM erc_20_collection_assets
WHERE
    owner = $1
LIMIT $2 OFFSET $3
`

type GetPaginated20AssetByOwnerAddressParams struct {
	Owner  string `json:"owner"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

func (q *Queries) GetPaginated20AssetByOwnerAddress(ctx context.Context, arg GetPaginated20AssetByOwnerAddressParams) ([]Erc20CollectionAsset, error) {
	rows, err := q.db.QueryContext(ctx, getPaginated20AssetByOwnerAddress, arg.Owner, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Erc20CollectionAsset
	for rows.Next() {
		var i Erc20CollectionAsset
		if err := rows.Scan(
			&i.ID,
			&i.ChainID,
			&i.AssetID,
			&i.Owner,
			&i.Balance,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const update20Asset = `-- name: Update20Asset :exec
UPDATE erc_20_collection_assets
SET
    owner = $2 
WHERE 
    id = $1
`

type Update20AssetParams struct {
	ID    uuid.UUID `json:"id"`
	Owner string    `json:"owner"`
}

func (q *Queries) Update20Asset(ctx context.Context, arg Update20AssetParams) error {
	_, err := q.db.ExecContext(ctx, update20Asset, arg.ID, arg.Owner)
	return err
}

// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: erc_1155_asset.sql

package db

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
)

const add1155Asset = `-- name: Add1155Asset :exec
INSERT INTO
    erc_1155_collection_assets (asset_id, chain_id, token_id, owner, balance, attributes)
VALUES (
    $1, $2, $3, $4, $5, $6
) ON CONFLICT ON CONSTRAINT UC_ERC1155 DO UPDATE SET
    balance = $5,
    attributes = $6
    
RETURNING id, chain_id, asset_id, token_id, owner, balance, attributes, created_at, updated_at
`

type Add1155AssetParams struct {
	AssetID    string         `json:"assetId"`
	ChainID    int32          `json:"chainId"`
	TokenID    string         `json:"tokenId"`
	Owner      string         `json:"owner"`
	Balance    string         `json:"balance"`
	Attributes sql.NullString `json:"attributes"`
}

func (q *Queries) Add1155Asset(ctx context.Context, arg Add1155AssetParams) error {
	_, err := q.db.ExecContext(ctx, add1155Asset,
		arg.AssetID,
		arg.ChainID,
		arg.TokenID,
		arg.Owner,
		arg.Balance,
		arg.Attributes,
	)
	return err
}

const count1155AssetByAssetId = `-- name: Count1155AssetByAssetId :one
SELECT COUNT(*) FROM erc_1155_collection_assets 
WHERE asset_id = $1
`

func (q *Queries) Count1155AssetByAssetId(ctx context.Context, assetID string) (int64, error) {
	row := q.db.QueryRowContext(ctx, count1155AssetByAssetId, assetID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const count1155AssetByOwner = `-- name: Count1155AssetByOwner :one
SELECT COUNT(*) FROM erc_1155_collection_assets 
WHERE owner = $1
`

func (q *Queries) Count1155AssetByOwner(ctx context.Context, owner string) (int64, error) {
	row := q.db.QueryRowContext(ctx, count1155AssetByOwner, owner)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const delete1155Asset = `-- name: Delete1155Asset :exec
DELETE 
FROM erc_1155_collection_assets
WHERE
    id = $1
`

func (q *Queries) Delete1155Asset(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, delete1155Asset, id)
	return err
}

const get1155AssetByAssetIdAndTokenId = `-- name: Get1155AssetByAssetIdAndTokenId :one
SELECT id, chain_id, asset_id, token_id, owner, balance, attributes, created_at, updated_at FROM erc_1155_collection_assets
WHERE
    asset_id = $1
    AND token_id = $2
`

type Get1155AssetByAssetIdAndTokenIdParams struct {
	AssetID string `json:"assetId"`
	TokenID string `json:"tokenId"`
}

func (q *Queries) Get1155AssetByAssetIdAndTokenId(ctx context.Context, arg Get1155AssetByAssetIdAndTokenIdParams) (Erc1155CollectionAsset, error) {
	row := q.db.QueryRowContext(ctx, get1155AssetByAssetIdAndTokenId, arg.AssetID, arg.TokenID)
	var i Erc1155CollectionAsset
	err := row.Scan(
		&i.ID,
		&i.ChainID,
		&i.AssetID,
		&i.TokenID,
		&i.Owner,
		&i.Balance,
		&i.Attributes,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getDetailERC1155Assets = `-- name: GetDetailERC1155Assets :one
SELECT 
    ts.asset_id,
    ts.token_id,
    ts.total_supply AS total_supply,
    json_agg(ca) AS asset_owners
FROM 
    erc_1155_total_supply ts
JOIN 
    erc_1155_collection_assets ca 
ON 
    ts.asset_id = ca.asset_id AND ts.token_id = ca.token_id
WHERE ts.asset_id = $1
AND ts.token_id = $2
GROUP BY 
    ts.asset_id, ts.token_id, ts.total_supply
`

type GetDetailERC1155AssetsParams struct {
	AssetID string `json:"assetId"`
	TokenID string `json:"tokenId"`
}

type GetDetailERC1155AssetsRow struct {
	AssetID     string          `json:"assetId"`
	TokenID     string          `json:"tokenId"`
	TotalSupply int64           `json:"totalSupply"`
	AssetOwners json.RawMessage `json:"assetOwners"`
}

func (q *Queries) GetDetailERC1155Assets(ctx context.Context, arg GetDetailERC1155AssetsParams) (GetDetailERC1155AssetsRow, error) {
	row := q.db.QueryRowContext(ctx, getDetailERC1155Assets, arg.AssetID, arg.TokenID)
	var i GetDetailERC1155AssetsRow
	err := row.Scan(
		&i.AssetID,
		&i.TokenID,
		&i.TotalSupply,
		&i.AssetOwners,
	)
	return i, err
}

const getPaginated1155AssetByAssetId = `-- name: GetPaginated1155AssetByAssetId :many
SELECT id, chain_id, asset_id, token_id, owner, balance, attributes, created_at, updated_at FROM erc_1155_collection_assets 
WHERE asset_id = $1
LIMIT $2 OFFSET $3
`

type GetPaginated1155AssetByAssetIdParams struct {
	AssetID string `json:"assetId"`
	Limit   int32  `json:"limit"`
	Offset  int32  `json:"offset"`
}

func (q *Queries) GetPaginated1155AssetByAssetId(ctx context.Context, arg GetPaginated1155AssetByAssetIdParams) ([]Erc1155CollectionAsset, error) {
	rows, err := q.db.QueryContext(ctx, getPaginated1155AssetByAssetId, arg.AssetID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Erc1155CollectionAsset
	for rows.Next() {
		var i Erc1155CollectionAsset
		if err := rows.Scan(
			&i.ID,
			&i.ChainID,
			&i.AssetID,
			&i.TokenID,
			&i.Owner,
			&i.Balance,
			&i.Attributes,
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

const getPaginated1155AssetByOwnerAddress = `-- name: GetPaginated1155AssetByOwnerAddress :many
SELECT id, chain_id, asset_id, token_id, owner, balance, attributes, created_at, updated_at FROM erc_1155_collection_assets
WHERE
    owner = $1
LIMIT $2 OFFSET $3
`

type GetPaginated1155AssetByOwnerAddressParams struct {
	Owner  string `json:"owner"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

func (q *Queries) GetPaginated1155AssetByOwnerAddress(ctx context.Context, arg GetPaginated1155AssetByOwnerAddressParams) ([]Erc1155CollectionAsset, error) {
	rows, err := q.db.QueryContext(ctx, getPaginated1155AssetByOwnerAddress, arg.Owner, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Erc1155CollectionAsset
	for rows.Next() {
		var i Erc1155CollectionAsset
		if err := rows.Scan(
			&i.ID,
			&i.ChainID,
			&i.AssetID,
			&i.TokenID,
			&i.Owner,
			&i.Balance,
			&i.Attributes,
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

const update1155Asset = `-- name: Update1155Asset :exec
UPDATE erc_1155_collection_assets
SET
    owner = $2 
WHERE 
    id = $1
`

type Update1155AssetParams struct {
	ID    uuid.UUID `json:"id"`
	Owner string    `json:"owner"`
}

func (q *Queries) Update1155Asset(ctx context.Context, arg Update1155AssetParams) error {
	_, err := q.db.ExecContext(ctx, update1155Asset, arg.ID, arg.Owner)
	return err
}

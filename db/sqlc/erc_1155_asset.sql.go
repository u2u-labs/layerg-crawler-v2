// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: erc_1155_asset.sql

package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

const add1155Asset = `-- name: Add1155Asset :exec
INSERT INTO
    erc_1155_collection_assets (asset_id, token_id, owner, balance, attributes)
VALUES (
    $1, $2, $3, $4, $5
) RETURNING id, chain_id, asset_id, token_id, owner, balance, attributes, created_at, updated_at
`

type Add1155AssetParams struct {
	AssetID    string
	TokenID    int32
	Owner      string
	Balance    int32
	Attributes pqtype.NullRawMessage
}

func (q *Queries) Add1155Asset(ctx context.Context, arg Add1155AssetParams) error {
	_, err := q.db.ExecContext(ctx, add1155Asset,
		arg.AssetID,
		arg.TokenID,
		arg.Owner,
		arg.Balance,
		arg.Attributes,
	)
	return err
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

const get1155AssetByAssetId = `-- name: Get1155AssetByAssetId :many
SELECT id, chain_id, asset_id, token_id, owner, balance, attributes, created_at, updated_at FROM erc_1155_collection_assets WHERE asset_id = $1
`

func (q *Queries) Get1155AssetByAssetId(ctx context.Context, assetID string) ([]Erc1155CollectionAsset, error) {
	rows, err := q.db.QueryContext(ctx, get1155AssetByAssetId, assetID)
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

const get1155AssetByAssetIdAndTokenId = `-- name: Get1155AssetByAssetIdAndTokenId :one
SELECT id, chain_id, asset_id, token_id, owner, balance, attributes, created_at, updated_at FROM erc_1155_collection_assets
WHERE
    asset_id = $1
    AND token_id = $2
`

type Get1155AssetByAssetIdAndTokenIdParams struct {
	AssetID string
	TokenID int32
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

const get1155AssetByOwner = `-- name: Get1155AssetByOwner :many
SELECT id, chain_id, asset_id, token_id, owner, balance, attributes, created_at, updated_at FROM erc_1155_collection_assets
WHERE
    owner = $1
`

func (q *Queries) Get1155AssetByOwner(ctx context.Context, owner string) ([]Erc1155CollectionAsset, error) {
	rows, err := q.db.QueryContext(ctx, get1155AssetByOwner, owner)
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
	ID    uuid.UUID
	Owner string
}

func (q *Queries) Update1155Asset(ctx context.Context, arg Update1155AssetParams) error {
	_, err := q.db.ExecContext(ctx, update1155Asset, arg.ID, arg.Owner)
	return err
}

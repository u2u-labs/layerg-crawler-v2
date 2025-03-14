// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: custom.sql

package graphqldb

import (
	"context"
)

const getItemByTokenId = `-- name: GetItemByTokenId :one
SELECT id, token_id, token_uri, standard, created_at FROM "item" WHERE token_id = $1
`

func (q *Queries) GetItemByTokenId(ctx context.Context, tokenID string) (Item, error) {
	row := q.db.QueryRowContext(ctx, getItemByTokenId, tokenID)
	var i Item
	err := row.Scan(
		&i.ID,
		&i.TokenID,
		&i.TokenUri,
		&i.Standard,
		&i.CreatedAt,
	)
	return i, err
}

const getOrCreateUser = `-- name: GetOrCreateUser :one
INSERT INTO "user" ("id") 
VALUES ($1) 
ON CONFLICT (id) DO UPDATE SET id = EXCLUDED.id
RETURNING id, created_at
`

// Add a new query to get or create user
func (q *Queries) GetOrCreateUser(ctx context.Context, id string) (User, error) {
	row := q.db.QueryRowContext(ctx, getOrCreateUser, id)
	var i User
	err := row.Scan(&i.ID, &i.CreatedAt)
	return i, err
}

const getUserBalance = `-- name: GetUserBalance :one
SELECT id, value FROM balance 
WHERE owner_id = $1 AND item_id = $2
LIMIT 1
`

type GetUserBalanceParams struct {
	OwnerID string `json:"owner_id"`
	ItemID  string `json:"item_id"`
}

type GetUserBalanceRow struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

func (q *Queries) GetUserBalance(ctx context.Context, arg GetUserBalanceParams) (GetUserBalanceRow, error) {
	row := q.db.QueryRowContext(ctx, getUserBalance, arg.OwnerID, arg.ItemID)
	var i GetUserBalanceRow
	err := row.Scan(&i.ID, &i.Value)
	return i, err
}

const upsertBalance = `-- name: UpsertBalance :one
INSERT INTO balance (
    id,
    item_id,
    owner_id,
    value,
    updated_at,
    contract
)
VALUES (
    $1, -- id (UUID)
    $2, -- item_id
    $3, -- owner_id
    $4, -- value
    $5, -- updated_at (block timestamp)
    $6  -- contract address
)
ON CONFLICT (id) 
DO UPDATE SET 
    value = EXCLUDED.value,
    updated_at = EXCLUDED.updated_at
RETURNING id, item_id, owner_id, value, updated_at, contract, created_at
`

type UpsertBalanceParams struct {
	ID        string `json:"id"`
	ItemID    string `json:"item_id"`
	OwnerID   string `json:"owner_id"`
	Value     string `json:"value"`
	UpdatedAt string `json:"updated_at"`
	Contract  string `json:"contract"`
}

func (q *Queries) UpsertBalance(ctx context.Context, arg UpsertBalanceParams) (Balance, error) {
	row := q.db.QueryRowContext(ctx, upsertBalance,
		arg.ID,
		arg.ItemID,
		arg.OwnerID,
		arg.Value,
		arg.UpdatedAt,
		arg.Contract,
	)
	var i Balance
	err := row.Scan(
		&i.ID,
		&i.ItemID,
		&i.OwnerID,
		&i.Value,
		&i.UpdatedAt,
		&i.Contract,
		&i.CreatedAt,
	)
	return i, err
}

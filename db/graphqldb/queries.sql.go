// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: queries.sql

package graphqldb

import (
	"context"
)

const createItem = `-- name: CreateItem :one
INSERT INTO "item" ("id", "token_id", "token_uri", "owner_id") VALUES ($1, $2, $3, $4) RETURNING id, token_id, token_uri, owner_id, created_at, user_id
`

type CreateItemParams struct {
	ID       string `json:"id"`
	TokenID  string `json:"token_id"`
	TokenUri string `json:"token_uri"`
	OwnerID  string `json:"owner_id"`
}

func (q *Queries) CreateItem(ctx context.Context, arg CreateItemParams) (Item, error) {
	row := q.db.QueryRowContext(ctx, createItem,
		arg.ID,
		arg.TokenID,
		arg.TokenUri,
		arg.OwnerID,
	)
	var i Item
	err := row.Scan(
		&i.ID,
		&i.TokenID,
		&i.TokenUri,
		&i.OwnerID,
		&i.CreatedAt,
		&i.UserID,
	)
	return i, err
}

const createMetadataUpdateRecord = `-- name: CreateMetadataUpdateRecord :one
INSERT INTO "metadata_update_record" ("id", "token_id", "actor_id") VALUES ($1, $2, $3) RETURNING id, token_id, actor_id, created_at
`

type CreateMetadataUpdateRecordParams struct {
	ID      string `json:"id"`
	TokenID string `json:"token_id"`
	ActorID string `json:"actor_id"`
}

func (q *Queries) CreateMetadataUpdateRecord(ctx context.Context, arg CreateMetadataUpdateRecordParams) (MetadataUpdateRecord, error) {
	row := q.db.QueryRowContext(ctx, createMetadataUpdateRecord, arg.ID, arg.TokenID, arg.ActorID)
	var i MetadataUpdateRecord
	err := row.Scan(
		&i.ID,
		&i.TokenID,
		&i.ActorID,
		&i.CreatedAt,
	)
	return i, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO "user" ("id") VALUES ($1) RETURNING id, created_at
`

func (q *Queries) CreateUser(ctx context.Context, id string) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, id)
	var i User
	err := row.Scan(&i.ID, &i.CreatedAt)
	return i, err
}

const deleteItem = `-- name: DeleteItem :exec
DELETE FROM "item" WHERE id = $1
`

func (q *Queries) DeleteItem(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, deleteItem, id)
	return err
}

const deleteMetadataUpdateRecord = `-- name: DeleteMetadataUpdateRecord :exec
DELETE FROM "metadata_update_record" WHERE id = $1
`

func (q *Queries) DeleteMetadataUpdateRecord(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, deleteMetadataUpdateRecord, id)
	return err
}

const getItem = `-- name: GetItem :one
SELECT id, token_id, token_uri, owner_id, created_at, user_id FROM "item" WHERE id = $1
`

func (q *Queries) GetItem(ctx context.Context, id string) (Item, error) {
	row := q.db.QueryRowContext(ctx, getItem, id)
	var i Item
	err := row.Scan(
		&i.ID,
		&i.TokenID,
		&i.TokenUri,
		&i.OwnerID,
		&i.CreatedAt,
		&i.UserID,
	)
	return i, err
}

const getMetadataUpdateRecord = `-- name: GetMetadataUpdateRecord :one
SELECT id, token_id, actor_id, created_at FROM "metadata_update_record" WHERE id = $1
`

func (q *Queries) GetMetadataUpdateRecord(ctx context.Context, id string) (MetadataUpdateRecord, error) {
	row := q.db.QueryRowContext(ctx, getMetadataUpdateRecord, id)
	var i MetadataUpdateRecord
	err := row.Scan(
		&i.ID,
		&i.TokenID,
		&i.ActorID,
		&i.CreatedAt,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id, created_at FROM "user" WHERE id = $1
`

func (q *Queries) GetUser(ctx context.Context, id string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(&i.ID, &i.CreatedAt)
	return i, err
}

const listItem = `-- name: ListItem :many
SELECT id, token_id, token_uri, owner_id, created_at, user_id FROM "item"
`

func (q *Queries) ListItem(ctx context.Context) ([]Item, error) {
	rows, err := q.db.QueryContext(ctx, listItem)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Item{}
	for rows.Next() {
		var i Item
		if err := rows.Scan(
			&i.ID,
			&i.TokenID,
			&i.TokenUri,
			&i.OwnerID,
			&i.CreatedAt,
			&i.UserID,
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

const listMetadataUpdateRecord = `-- name: ListMetadataUpdateRecord :many
SELECT id, token_id, actor_id, created_at FROM "metadata_update_record"
`

func (q *Queries) ListMetadataUpdateRecord(ctx context.Context) ([]MetadataUpdateRecord, error) {
	rows, err := q.db.QueryContext(ctx, listMetadataUpdateRecord)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []MetadataUpdateRecord{}
	for rows.Next() {
		var i MetadataUpdateRecord
		if err := rows.Scan(
			&i.ID,
			&i.TokenID,
			&i.ActorID,
			&i.CreatedAt,
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

const listUser = `-- name: ListUser :many
SELECT id, created_at FROM "user"
`

func (q *Queries) ListUser(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, listUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(&i.ID, &i.CreatedAt); err != nil {
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

const updateItem = `-- name: UpdateItem :one
UPDATE "item" SET "token_id" = $2, "token_uri" = $3, "owner_id" = $4 WHERE id = $1 RETURNING id, token_id, token_uri, owner_id, created_at, user_id
`

type UpdateItemParams struct {
	ID       string `json:"id"`
	TokenID  string `json:"token_id"`
	TokenUri string `json:"token_uri"`
	OwnerID  string `json:"owner_id"`
}

func (q *Queries) UpdateItem(ctx context.Context, arg UpdateItemParams) (Item, error) {
	row := q.db.QueryRowContext(ctx, updateItem,
		arg.ID,
		arg.TokenID,
		arg.TokenUri,
		arg.OwnerID,
	)
	var i Item
	err := row.Scan(
		&i.ID,
		&i.TokenID,
		&i.TokenUri,
		&i.OwnerID,
		&i.CreatedAt,
		&i.UserID,
	)
	return i, err
}

const updateMetadataUpdateRecord = `-- name: UpdateMetadataUpdateRecord :one
UPDATE "metadata_update_record" SET "token_id" = $2, "actor_id" = $3 WHERE id = $1 RETURNING id, token_id, actor_id, created_at
`

type UpdateMetadataUpdateRecordParams struct {
	ID      string `json:"id"`
	TokenID string `json:"token_id"`
	ActorID string `json:"actor_id"`
}

func (q *Queries) UpdateMetadataUpdateRecord(ctx context.Context, arg UpdateMetadataUpdateRecordParams) (MetadataUpdateRecord, error) {
	row := q.db.QueryRowContext(ctx, updateMetadataUpdateRecord, arg.ID, arg.TokenID, arg.ActorID)
	var i MetadataUpdateRecord
	err := row.Scan(
		&i.ID,
		&i.TokenID,
		&i.ActorID,
		&i.CreatedAt,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :exec

DELETE FROM "user" WHERE id = $1
`

// Skip update query generation as there are no updateable fields
func (q *Queries) UpdateUser(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, updateUser, id)
	return err
}

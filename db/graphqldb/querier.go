// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package graphqldb

import (
	"context"
)

type Querier interface {
	CreateValue(ctx context.Context, arg CreateValueParams) (Value, error)
	DeleteValue(ctx context.Context, id string) error
	GetValue(ctx context.Context, id string) (Value, error)
	ListValue(ctx context.Context) ([]Value, error)
	UpdateValue(ctx context.Context, arg UpdateValueParams) (Value, error)
}

var _ Querier = (*Queries)(nil)

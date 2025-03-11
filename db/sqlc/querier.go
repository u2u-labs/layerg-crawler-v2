// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"context"
)

type Querier interface {
	CreateAsset(ctx context.Context, arg CreateAssetParams) (Asset, error)
	CreateChain(ctx context.Context, arg CreateChainParams) (Chain, error)
	CreateOnchainHistory(ctx context.Context, arg CreateOnchainHistoryParams) (OnchainHistory, error)
	GetAllChain(ctx context.Context) ([]Chain, error)
	GetAssetByAddress(ctx context.Context, arg GetAssetByAddressParams) (Asset, error)
	GetChainById(ctx context.Context, id int32) (Chain, error)
	GetOnchainHistoriesByAsset(ctx context.Context, arg GetOnchainHistoriesByAssetParams) ([]OnchainHistory, error)
	GetPaginatedAssetsByChainId(ctx context.Context, arg GetPaginatedAssetsByChainIdParams) ([]Asset, error)
	UpdateChainLatestBlock(ctx context.Context, arg UpdateChainLatestBlockParams) error
}

var _ Querier = (*Queries)(nil)

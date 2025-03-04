// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Asset struct {
	ID              string        `json:"id"`
	ChainID         int32         `json:"chain_id"`
	ContractAddress string        `json:"contract_address"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
	InitialBlock    sql.NullInt64 `json:"initial_block"`
	LastUpdated     sql.NullTime  `json:"last_updated"`
}

type Chain struct {
	ID          int32  `json:"id"`
	Chain       string `json:"chain"`
	Name        string `json:"name"`
	RpcUrl      string `json:"rpc_url"`
	ChainID     int64  `json:"chain_id"`
	Explorer    string `json:"explorer"`
	LatestBlock int64  `json:"latest_block"`
	BlockTime   int32  `json:"block_time"`
}

type OnchainHistory struct {
	ID        uuid.UUID       `json:"id"`
	From      string          `json:"from"`
	To        string          `json:"to"`
	ChainID   int32           `json:"chain_id"`
	AssetID   string          `json:"asset_id"`
	TxHash    string          `json:"tx_hash"`
	Receipt   json.RawMessage `json:"receipt"`
	EventType sql.NullString  `json:"event_type"`
	Timestamp time.Time       `json:"timestamp"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

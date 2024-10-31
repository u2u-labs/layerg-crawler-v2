// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

type AssetType string

const (
	AssetTypeERC721  AssetType = "ERC721"
	AssetTypeERC1155 AssetType = "ERC1155"
	AssetTypeERC20   AssetType = "ERC20"
)

func (e *AssetType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = AssetType(s)
	case string:
		*e = AssetType(s)
	default:
		return fmt.Errorf("unsupported scan type for AssetType: %T", src)
	}
	return nil
}

type NullAssetType struct {
	AssetType AssetType `json:"assetType"`
	Valid     bool      `json:"valid"` // Valid is true if AssetType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullAssetType) Scan(value interface{}) error {
	if value == nil {
		ns.AssetType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.AssetType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullAssetType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.AssetType), nil
}

type App struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	SecretKey string    `json:"secretKey"`
}

type Asset struct {
	ID                string        `json:"id"`
	ChainID           int32         `json:"chainId"`
	CollectionAddress string        `json:"collectionAddress"`
	Type              AssetType     `json:"type"`
	CreatedAt         time.Time     `json:"createdAt"`
	UpdatedAt         time.Time     `json:"updatedAt"`
	DecimalData       sql.NullInt16 `json:"decimalData"`
	InitialBlock      sql.NullInt64 `json:"initialBlock"`
	LastUpdated       sql.NullTime  `json:"lastUpdated"`
}

type Chain struct {
	ID          int32  `json:"id"`
	Chain       string `json:"chain"`
	Name        string `json:"name"`
	RpcUrl      string `json:"rpcUrl"`
	ChainID     int64  `json:"chainId"`
	Explorer    string `json:"explorer"`
	LatestBlock int64  `json:"latestBlock"`
	BlockTime   int32  `json:"blockTime"`
}

type Erc1155CollectionAsset struct {
	ID         uuid.UUID             `json:"id"`
	ChainID    int32                 `json:"chainId"`
	AssetID    string                `json:"assetId"`
	TokenID    string                `json:"tokenId"`
	Owner      string                `json:"owner"`
	Balance    string                `json:"balance"`
	Attributes pqtype.NullRawMessage `json:"attributes"`
	CreatedAt  time.Time             `json:"createdAt"`
	UpdatedAt  time.Time             `json:"updatedAt"`
}

type Erc20CollectionAsset struct {
	ID        uuid.UUID `json:"id"`
	ChainID   int32     `json:"chainId"`
	AssetID   string    `json:"assetId"`
	Owner     string    `json:"owner"`
	Balance   string    `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Erc721CollectionAsset struct {
	ID         uuid.UUID             `json:"id"`
	ChainID    int32                 `json:"chainId"`
	AssetID    string                `json:"assetId"`
	TokenID    string                `json:"tokenId"`
	Owner      string                `json:"owner"`
	Attributes pqtype.NullRawMessage `json:"attributes"`
	CreatedAt  time.Time             `json:"createdAt"`
	UpdatedAt  time.Time             `json:"updatedAt"`
}

type OnchainHistory struct {
	ID        uuid.UUID `json:"id"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	AssetID   string    `json:"assetId"`
	TokenID   string    `json:"tokenId"`
	Amount    float64   `json:"amount"`
	TxHash    string    `json:"txHash"`
	Timestamp time.Time `json:"timestamp"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

package types

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Metadata struct {
	Decimal      uint8
	InitialBlock *pgtype.Numeric
	LastUpdated  time.Time
}

type Asset struct {
	ChainId           int
	CollectionAddress string
	Type              string
	Metadata          Metadata
}

type ERC721CollectionAsset struct {
	TokenId *pgtype.Numeric `gorm:"index:,unique,sort:desc"`
	Owner   string          `gorm:"index:,type:btree,length:20"`
	URI     string
}

type ERC20Asset struct {
	Owner        string `gorm:"index:,unique,sort:desc,length:20"`
	Balance      string
	BalanceFloat float64 `gorm:"sort:desc"`
}

type ERC1155CollectionAsset struct {
	TokenId *pgtype.Numeric `gorm:"type:numeric;index:idx_contract_token,unique,sort:desc"`
	Owner   string          `gorm:"index:,type:btree,length:20"`
	Balance *pgtype.Numeric
}

type OnchainHistory struct {
	From              string `gorm:"index:,type:btree,length:20"`
	To                string `gorm:"index:,type:btree,length:20"`
	CollectionAddress string `gorm:"index:,type:btree,length:20"`
	TokenId           *pgtype.Numeric
	Amount            *pgtype.Numeric
	TxHash            string    `gorm:"index:,type:btree,length:32"`
	Timestamp         time.Time `gorm:"sort:desc"`
}

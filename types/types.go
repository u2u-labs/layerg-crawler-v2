package types

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"time"

	"github.com/unicornultrafoundation/go-u2u/common"
)

type Network struct {
	Chain    string `yaml:"chain"`
	Name     string `yaml:"name"`
	RpcUrl   string `yaml:"rpcUrl"`
	ChainId  int    `yaml:"chainId"`
	Explorer string `yaml:"explorer"`
}

func (n *Network) String() string {
	return fmt.Sprintf(`%s %s
RPC URL: %s
Chain ID: %d,
Explorer: %s`, n.Chain, n.Name, n.RpcUrl, n.ChainId, n.Explorer)
}

type Metadata struct {
	Decimal      uint8
	InitialBlock pgtype.Numeric
	LastUpdated  time.Time
}

type Asset struct {
	Chain             string
	CollectionAddress *common.Address
	Type              string
	Metadata          Metadata
}

type ERC721CollectionAsset struct {
	TokenId pgtype.Numeric  `gorm:"index:,unique,sort:desc"`
	Owner   *common.Address `gorm:"index:,type:btree,length:20"`
}

type ERC20Asset struct {
	Owner        *common.Address `gorm:"index:,unique,sort:desc,length:20"`
	Balance      string
	BalanceFloat float64 `gorm:"sort:desc"`
}

type ERC1155CollectionAsset struct {
	TokenId pgtype.Numeric  `gorm:"type:numeric;index:idx_contract_token,unique,sort:desc"`
	Owner   *common.Address `gorm:"index:,type:btree,length:20"`
	Balance pgtype.Numeric
}

type OnchainHistory struct {
	From              *common.Address `gorm:"index:,type:btree,length:20"`
	To                *common.Address `gorm:"index:,type:btree,length:20"`
	CollectionAddress *common.Address `gorm:"index:,type:btree,length:20"`
	TokenId           pgtype.Numeric
	Amount            pgtype.Numeric
	TxHash            *common.Hash `gorm:"index:,type:btree,length:32"`
	Timestamp         time.Time    `gorm:"sort:desc"`
}

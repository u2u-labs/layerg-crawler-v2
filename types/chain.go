package types

import (
	"fmt"
)

type Chain struct {
	Id          int    `json:"id" gorm:"index:,unique,sort:asc"`
	Chain       string `json:"chain" yaml:"chain"`
	Name        string `json:"name" yaml:"name"`
	RpcUrl      string `json:"rpcUrl" yaml:"rpcUrl"`
	ChainId     int    `json:"chainId" yaml:"chainId"`
	Explorer    string `json:"explorer" yaml:"explorer"`
	LatestBlock uint64 `json:"latestBlock" `
	BlockTime   uint   `json:"blockTime"`
}

type Contract struct {
	Id          int    `json:"id" gorm:"index:,unique,sort:asc"`
	Address     string `json:"address" gorm:"index:,unique,sort:asc"`
	LatestBlock uint64
}

func (n *Chain) String() string {
	return fmt.Sprintf(`%s %s - RPC URL: %s - Chain ID: %d - Explorer: %s - Latest Block: %d`,
		n.Chain, n.Name, n.RpcUrl, n.ChainId, n.Explorer, n.LatestBlock)
}

package types

import (
	"fmt"
)

type Chain struct {
	Id             int    `json:"id" gorm:"index:,unique,sort:asc"`
	Chain          string `json:"chain" yaml:"chain"`
	Name           string `json:"name" yaml:"name"`
	RpcUrl         string `json:"rpcUrl" yaml:"rpcUrl"`
	ChainId        int    `json:"chainId" yaml:"chainId"`
	Explorer       string `json:"explorer" yaml:"explorer"`
	TokenContracts string `json:"tokenContracts"`
	NftContracts   string `json:"nftContracts"`
	LatestBlock    uint64 `json:"latestBlock" `
	BlockTime      uint   `json:"blockTime"`
}

func (n *Chain) String() string {
	return fmt.Sprintf(`%s %s - RPC URL: %s - Chain ID: %d - Explorer: %s - Token Contracts: %v - NFT Contracts: %v - Latest Block: %d`,
		n.Chain, n.Name, n.RpcUrl, n.ChainId, n.Explorer, n.TokenContracts, n.NftContracts, n.LatestBlock)
}

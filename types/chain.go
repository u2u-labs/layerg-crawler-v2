package types

import "fmt"

type Network struct {
	Id          int    `json:"id" gorm:"index:,unique,sort:asc"`
	Chain       string `json:"chain" yaml:"chain"`
	Name        string `json:"name" yaml:"name"`
	RpcUrl      string `json:"rpcUrl" yaml:"rpcUrl"`
	ChainId     int    `json:"chainId" yaml:"chainId"`
	Explorer    string `json:"explorer" yaml:"explorer"`
	LatestBlock uint64 `json:"latestBlock" `
}

func (n *Network) String() string {
	return fmt.Sprintf(`%s %s
RPC URL: %s
Chain ID: %d,
Explorer: %s
LatesBlock: %d`, n.Chain, n.Name, n.RpcUrl, n.ChainId, n.Explorer, n.LatestBlock)
}

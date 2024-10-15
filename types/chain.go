package types

import "fmt"

type Network struct {
	Id          int    `gorm:"index:,unique,sort:asc"`
	Chain       string `yaml:"chain"`
	Name        string `yaml:"name"`
	RpcUrl      string `yaml:"rpcUrl"`
	ChainId     int    `yaml:"chainId"`
	Explorer    string `yaml:"explorer"`
	LatestBlock uint64
}

func (n *Network) String() string {
	return fmt.Sprintf(`%s %s
RPC URL: %s
Chain ID: %d,
Explorer: %s
LatesBlock: %d`, n.Chain, n.Name, n.RpcUrl, n.ChainId, n.Explorer, n.LatestBlock)
}

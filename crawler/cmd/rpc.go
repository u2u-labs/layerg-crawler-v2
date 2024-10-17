package cmd

import (
	"github.com/unicornultrafoundation/go-u2u/ethclient"

	"layerg-crawler/types"
)

func initChainClient(chain *types.Network) (*ethclient.Client, error) {
	return ethclient.Dial(chain.RpcUrl)
}

package cmd

import (
	"github.com/unicornultrafoundation/go-u2u/ethclient"

	"github.com/u2u-labs/layerg-crawler/types"
)

func initChainClient(chain *types.Chain) (*ethclient.Client, error) {
	return ethclient.Dial(chain.RpcUrl)
}

package cmd

import (
	"github.com/unicornultrafoundation/go-u2u/ethclient"

	db "github.com/u2u-labs/layerg-crawler/db/sqlc"
)

func initChainClient(chain *db.Chain) (*ethclient.Client, error) {
	return ethclient.Dial(chain.RpcUrl)
}

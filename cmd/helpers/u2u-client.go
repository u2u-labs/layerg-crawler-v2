package helpers

import (
	"context"

	db "github.com/u2u-labs/layerg-crawler/db/sqlc"
	"github.com/unicornultrafoundation/go-u2u/ethclient"
	"github.com/unicornultrafoundation/go-u2u/rpc"
)

func GetLastestBlockFromChainUrl(url string) (uint64, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return 0, err
	}
	defer client.Close()

	latest, err := client.BlockNumber(context.Background())
	if err != nil {
		return 0, err
	}

	return latest, nil
}

func InitChainClient(chain db.Chain) (*ethclient.Client, error) {
	return ethclient.Dial(chain.RpcUrl)
}

func InitNewRPCClient(url string) (*rpc.Client, error) {
	return rpc.Dial(url)
}

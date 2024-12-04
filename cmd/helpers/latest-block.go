package helpers

import (
	"context"

	"github.com/unicornultrafoundation/go-u2u/ethclient"
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

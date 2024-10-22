package cmd

import (
	"context"
	"go.uber.org/zap"
	"time"

	utypes "github.com/unicornultrafoundation/go-u2u/core/types"
	"github.com/unicornultrafoundation/go-u2u/ethclient"

	"github.com/u2u-labs/layerg-crawler/types"
)

func StartChainCrawler(sugar *zap.SugaredLogger, client *ethclient.Client, chain *types.Chain) {
	ctx := context.Background()
	timer := time.NewTimer(time.Duration(chain.BlockTime) * time.Millisecond)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			// Process new blocks
			ProcessLatestBlocks(ctx, sugar, client, chain)
			timer.Reset(time.Duration(chain.BlockTime) * time.Millisecond)
		}
	}
}

func ProcessLatestBlocks(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client, chain *types.Chain) error {
	latest, err := client.BlockNumber(ctx)
	sugar.Info(latest)
	if err != nil {
		sugar.Errorw("Failed to fetch latest blocks", "err", err, "chain", chain)
		return err
	}
	//var receipts []*utypes.Receipt
	//for i := chain.LatestBlock + 1; i <= latest; i++ {
	//	r, err := client.BlockReceipts(ctx, rpc.BlockNumberOrHashWithNumber(rpc.BlockNumber(latest)))
	//	if err != nil {
	//		sugar.Errorw("Failed to fetch latest block receipts", "err", err, "height", i, "chain", chain)
	//		return err
	//	}
	//	receipts = append(receipts, r...)
	//}

	return nil
}

func FilterEvents(chain types.Chain, receipts utypes.Receipts) {
}

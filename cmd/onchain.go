package cmd

import (
	"context"
	"time"

	utypes "github.com/unicornultrafoundation/go-u2u/core/types"
	"github.com/unicornultrafoundation/go-u2u/ethclient"
	"github.com/unicornultrafoundation/go-u2u/rpc"
	"go.uber.org/zap"

	db "github.com/u2u-labs/layerg-crawler/db/sqlc"
)

func StartChainCrawler(sugar *zap.SugaredLogger, client *ethclient.Client, chain *db.Chain) {
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

func ProcessLatestBlocks(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client, chain *db.Chain) error {
	latest, err := client.BlockNumber(ctx)
	sugar.Info(latest)
	if err != nil {
		sugar.Errorw("Failed to fetch latest blocks", "err", err, "chain", chain)
		return err
	}
	var receipts []*utypes.Receipt
	for i := chain.LatestBlock + 1; i <= int64(latest); i++ {
		r, err := client.BlockReceipts(ctx, rpc.BlockNumberOrHashWithNumber(rpc.BlockNumber(latest)))
		if err != nil {
			sugar.Errorw("Failed to fetch latest block receipts", "err", err, "height", i, "chain", chain)
			return err
		}
		receipts = append(receipts, r...)
	}
	for _, r := range receipts {
		sugar.Info("hash", r.TxHash)
	}
	return nil
}

func FilterEvents(chain db.Chain, receipts utypes.Receipts) {
}

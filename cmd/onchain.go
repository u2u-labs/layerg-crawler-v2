package cmd

import (
	"context"
	"github.com/sqlc-dev/pqtype"
	"time"

	"github.com/unicornultrafoundation/go-u2u/common"
	utypes "github.com/unicornultrafoundation/go-u2u/core/types"
	"github.com/unicornultrafoundation/go-u2u/ethclient"
	"github.com/unicornultrafoundation/go-u2u/rpc"
	"go.uber.org/zap"

	"github.com/u2u-labs/layerg-crawler/cmd/utils"
	db "github.com/u2u-labs/layerg-crawler/db/sqlc"
)

func StartChainCrawler(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client, q *db.Queries, chain *db.Chain) {
	sugar.Infow("Start chain crawler", "chain", chain.Chain+" "+chain.Name)
	timer := time.NewTimer(time.Duration(chain.BlockTime) * time.Millisecond)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			// Process new blocks
			ProcessLatestBlocks(ctx, sugar, client, q, chain)
			timer.Reset(time.Duration(chain.BlockTime) * time.Millisecond)
		}
	}
}

func ProcessLatestBlocks(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client, q *db.Queries, chain *db.Chain) error {
	latest, err := client.BlockNumber(ctx)
	if err != nil {
		sugar.Errorw("Failed to fetch latest blocks", "err", err, "chain", chain)
		return err
	}
	// Process each block between
	for i := chain.LatestBlock + 1; i <= int64(latest); i++ {
		if i%50 == 0 {
			sugar.Infow("Importing block receipts", "chain", chain.Chain+" "+chain.Name, "block", i)
		}
		receipts, err := client.BlockReceipts(ctx, rpc.BlockNumberOrHashWithNumber(rpc.BlockNumber(latest)))
		if err != nil {
			sugar.Errorw("Failed to fetch latest block receipts", "err", err, "height", i, "chain", chain)
			return err
		}
		if err = FilterEvents(ctx, sugar, q, chain, receipts); err != nil {
			sugar.Errorw("Failed to filter events", "err", err, "height", i, "chain", chain)
			return err
		}
	}
	// Update latest block processed
	chain.LatestBlock = int64(latest)
	if err = q.UpdateChainLatestBlock(ctx, db.UpdateChainLatestBlockParams{
		ID:          chain.ID,
		LatestBlock: int64(latest),
	}); err != nil {
		sugar.Errorw("Failed to update chain latest blocks in DB", "err", err, "chain", chain)
		return err
	}
	return nil
}

func FilterEvents(ctx context.Context, sugar *zap.SugaredLogger, q *db.Queries, chain *db.Chain, receipts utypes.Receipts) error {
	for _, r := range receipts {
		for _, l := range r.Logs {
			switch contractType[chain.ID][l.Address.Hex()].Type {
			case db.AssetTypeERC20:
				if err := handleErc20Transfer(ctx, sugar, q, chain, l); err != nil {
					sugar.Errorw("handleErc20Transfer", "err", err)
					return err
				}
			case db.AssetTypeERC721:
				if err := handleErc721Transfer(ctx, sugar, q, chain, l); err != nil {
					sugar.Errorw("handleErc721Transfer", "err", err)
					return err
				}
			case db.AssetTypeERC1155:
				if l.Topics[0].Hex() == utils.TransferSingleSig {
					if err := handleErc1155TransferSingle(ctx, sugar, q, chain, l); err != nil {
						sugar.Errorw("handleErc1155TransferSingle", "err", err)
						return err
					}
				}
				if l.Topics[0].Hex() == utils.TransferBatchSig {
					if err := handleErc1155TransferBatch(ctx, sugar, q, chain, l); err != nil {
						sugar.Errorw("handleErc1155TransferBatch", "err", err)
						return err
					}
				}
			default:
				continue
			}
		}
	}
	return nil
}

func handleErc20Transfer(ctx context.Context, sugar *zap.SugaredLogger, q *db.Queries, chain *db.Chain, l *utypes.Log) error {
	if l.Topics[0].Hex() != utils.TransferEventSig {
		return nil
	}
	// Unpack the log data
	var event utils.Erc20TransferEvent
	err := utils.ERC20ABI.UnpackIntoInterface(&event, "Transfer", l.Data)
	if err != nil {
		sugar.Fatalf("Failed to unpack log: %v", err)
		return err
	}
	// Decode the indexed fields manually
	event.From = common.BytesToAddress(l.Topics[1].Bytes())
	event.To = common.BytesToAddress(l.Topics[2].Bytes())
	amount, _ := event.Value.Float64()

	if err = q.AddOnchainTransaction(ctx, db.AddOnchainTransactionParams{
		From:      event.From.Hex(),
		To:        event.To.Hex(),
		AssetID:   contractType[chain.ID][l.Address.Hex()].ID,
		TokenID:   "0",
		Amount:    amount,
		TxHash:    l.TxHash.Hex(),
		Timestamp: time.Now(),
	}); err != nil {
		return err
	}
	// Update holders without balances
	if err = q.Add20Asset(ctx, db.Add20AssetParams{
		AssetID: contractType[chain.ID][l.Address.Hex()].ID,
		ChainID: chain.ID,
		Owner:   event.From.Hex(),
		Balance: "0",
	}); err != nil {
		return err
	}
	if err = q.Add20Asset(ctx, db.Add20AssetParams{
		AssetID: contractType[chain.ID][l.Address.Hex()].ID,
		ChainID: chain.ID,
		Owner:   event.To.Hex(),
		Balance: "0",
	}); err != nil {
		return err
	}

	return nil
}

func handleErc721Transfer(ctx context.Context, sugar *zap.SugaredLogger, q *db.Queries, chain *db.Chain, l *utypes.Log) error {
	if l.Topics[0].Hex() != utils.TransferEventSig {
		return nil
	}
	// Decode the indexed fields manually
	event := utils.Erc721TransferEvent{
		From:    common.BytesToAddress(l.Topics[1].Bytes()),
		To:      common.BytesToAddress(l.Topics[2].Bytes()),
		TokenID: l.Topics[3].Big(),
	}
	err := q.AddOnchainTransaction(ctx, db.AddOnchainTransactionParams{
		From:      event.From.Hex(),
		To:        event.To.Hex(),
		AssetID:   contractType[chain.ID][l.Address.Hex()].ID,
		TokenID:   event.TokenID.String(),
		Amount:    0,
		TxHash:    l.TxHash.Hex(),
		Timestamp: time.Now(),
	})
	if err != nil {
		return err
	}
	// Update NFT holder
	if err = q.Add721Asset(ctx, db.Add721AssetParams{
		AssetID:    contractType[chain.ID][l.Address.Hex()].ID,
		ChainID:    chain.ID,
		TokenID:    event.TokenID.String(),
		Owner:      event.From.Hex(),
		Attributes: pqtype.NullRawMessage{},
	}); err != nil {
		return err
	}

	return nil
}

func handleErc1155TransferBatch(ctx context.Context, sugar *zap.SugaredLogger, q *db.Queries, chain *db.Chain, l *utypes.Log) error {
	// Decode TransferBatch log
	var event utils.Erc1155TransferBatchEvent
	err := utils.ERC1155ABI.UnpackIntoInterface(&event, "TransferBatch", l.Data)
	if err != nil {
		sugar.Errorw("Failed to unpack TransferBatch log:", "err", err)
	}

	// Decode the indexed fields for TransferBatch
	event.Operator = common.BytesToAddress(l.Topics[1].Bytes())
	event.From = common.BytesToAddress(l.Topics[2].Bytes())
	event.To = common.BytesToAddress(l.Topics[3].Bytes())

	for i := range event.Ids {
		amount, _ := event.Values[i].Float64()
		if err = q.AddOnchainTransaction(ctx, db.AddOnchainTransactionParams{
			From:      event.From.Hex(),
			To:        event.To.Hex(),
			AssetID:   contractType[chain.ID][l.Address.Hex()].ID,
			TokenID:   event.Ids[i].String(),
			Amount:    amount,
			TxHash:    l.TxHash.Hex(),
			Timestamp: time.Now(),
		}); err != nil {
			return err
		}
		if err = q.Add1155Asset(ctx, db.Add1155AssetParams{
			AssetID:    contractType[chain.ID][l.Address.Hex()].ID,
			ChainID:    chain.ID,
			TokenID:    event.Ids[i].String(),
			Owner:      event.To.Hex(),
			Balance:    event.Values[i].String(),
			Attributes: pqtype.NullRawMessage{},
		}); err != nil {
			return err
		}
	}

	return nil
}

func handleErc1155TransferSingle(ctx context.Context, sugar *zap.SugaredLogger, q *db.Queries, chain *db.Chain, l *utypes.Log) error {
	// Decode TransferSingle log
	var event utils.Erc1155TransferSingleEvent
	err := utils.ERC1155ABI.UnpackIntoInterface(&event, "TransferSingle", l.Data)
	if err != nil {
		sugar.Errorw("Failed to unpack TransferSingle log:", "err", err)
	}

	// Decode the indexed fields for TransferSingle
	event.Operator = common.BytesToAddress(l.Topics[1].Bytes())
	event.From = common.BytesToAddress(l.Topics[2].Bytes())
	event.To = common.BytesToAddress(l.Topics[3].Bytes())

	amount, _ := event.Value.Float64()
	if err = q.AddOnchainTransaction(ctx, db.AddOnchainTransactionParams{
		From:      event.From.Hex(),
		To:        event.To.Hex(),
		AssetID:   contractType[chain.ID][l.Address.Hex()].ID,
		TokenID:   event.Id.String(),
		Amount:    amount,
		TxHash:    l.TxHash.Hex(),
		Timestamp: time.Now(),
	}); err != nil {
		return err
	}
	if err = q.Add1155Asset(ctx, db.Add1155AssetParams{
		AssetID:    contractType[chain.ID][l.Address.Hex()].ID,
		ChainID:    chain.ID,
		TokenID:    event.Id.String(),
		Owner:      event.To.Hex(),
		Balance:    event.Value.String(),
		Attributes: pqtype.NullRawMessage{},
	}); err != nil {
		return err
	}

	return nil
}

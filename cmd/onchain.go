package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	u2u "github.com/unicornultrafoundation/go-u2u"
	"github.com/unicornultrafoundation/go-u2u/common"
	utypes "github.com/unicornultrafoundation/go-u2u/core/types"
	"github.com/unicornultrafoundation/go-u2u/ethclient"
	"github.com/unicornultrafoundation/go-u2u/rpc"
	"go.uber.org/zap"

	"github.com/u2u-labs/layerg-crawler/cmd/utils"
	rdb "github.com/u2u-labs/layerg-crawler/db"
	db "github.com/u2u-labs/layerg-crawler/db/sqlc"
)

func StartChainCrawler(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client, q *db.Queries, chain *db.Chain, rdb *redis.Client) {
	sugar.Infow("Start chain crawler", "chain", chain)
	timer := time.NewTimer(time.Duration(chain.BlockTime) * time.Millisecond)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:

			// Process new blocks
			ProcessLatestBlocks(ctx, sugar, client, q, chain, rdb)
			timer.Reset(time.Duration(chain.BlockTime) * time.Millisecond)
		}
	}
}

func StartBackfillCrawler(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client, q *db.Queries, chain *db.Chain, bf *db.GetCrawlingBackfillCrawlerRow, rdb *redis.Client) error {
	sugar.Infow("Start backfill crawler", "collectionAddress", bf.CollectionAddress, "chain", chain)
	timer := time.NewTimer(time.Duration(100) * time.Millisecond)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			// scan 1000 blocks at once
			toScanBlock := bf.CurrentBlock + 1000
			if bf.InitialBlock.Valid && toScanBlock >= bf.InitialBlock.Int64 {
				toScanBlock = bf.InitialBlock.Int64
				bf.Status = db.CrawlerStatusCRAWLED
			}

			var transferEventSig []string

			if bf.Type == db.AssetTypeERC20 || bf.Type == db.AssetTypeERC721 {
				transferEventSig = []string{utils.TransferEventSig}
			} else if bf.Type == db.AssetTypeERC1155 {
				transferEventSig = []string{utils.TransferSingleSig, utils.TransferBatchSig}
			}

			// Initialize a two-dimensional slice of common.Hash
			topics := make([][]common.Hash, len(transferEventSig))

			// Populate the topics slice
			for i, sig := range transferEventSig {
				topics[i] = []common.Hash{common.HexToHash(sig)} // Convert each string to common.Hash
			}

			logs, _ := client.FilterLogs(ctx, u2u.FilterQuery{
				Topics:    topics,
				BlockHash: nil,
				FromBlock: big.NewInt(bf.CurrentBlock),
				ToBlock:   big.NewInt(toScanBlock),
				Addresses: []common.Address{common.HexToAddress(bf.CollectionAddress)},
			})

			for _, l := range logs {

				switch bf.Type {
				case db.AssetTypeERC20:
					if err := handleErc20Transfer(ctx, sugar, q, client, chain, rdb, &l, true); err != nil {
						sugar.Errorw("handleErc20Transfer", "err", err)
						return err
					}
				case db.AssetTypeERC721:
					if err := handleErc721Transfer(ctx, sugar, q, client, chain, rdb, &l, true); err != nil {
						sugar.Errorw("handleErc721Transfer", "err", err)
						return err
					}
				case db.AssetTypeERC1155:
					if l.Topics[0].Hex() == utils.TransferSingleSig {
						if err := handleErc1155TransferSingle(ctx, sugar, q, client, chain, rdb, &l, true); err != nil {
							sugar.Errorw("handleErc1155TransferSingle", "err", err)
							return err
						}
					}
					if l.Topics[0].Hex() == utils.TransferBatchSig {
						if err := handleErc1155TransferBatch(ctx, sugar, q, client, chain, rdb, &l, true); err != nil {
							sugar.Errorw("handleErc1155TransferBatch", "err", err)
							return err
						}
					}
				default:
					continue
				}
			}
			bf.CurrentBlock = toScanBlock

			q.UpdateCrawlingBackfill(ctx, db.UpdateCrawlingBackfillParams{
				ChainID:           bf.ChainID,
				CollectionAddress: bf.CollectionAddress,
				Status:            bf.Status,
				CurrentBlock:      bf.CurrentBlock,
			})

			if bf.Status == db.CrawlerStatusCRAWLED {
				return nil
			}

			timer.Reset(time.Duration(100) * time.Millisecond)
		}
	}
}

func ProcessLatestBlocks(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client, q *db.Queries, chain *db.Chain, rdb *redis.Client) error {
	latest, err := client.BlockNumber(ctx)
	if err != nil {
		sugar.Errorw("Failed to fetch latest blocks", "err", err, "chain", chain)
		return err
	}
	// Process each block between
	for i := chain.LatestBlock + 1; i <= int64(latest); i++ {
		if i%50 == 0 {
			sugar.Infow("Importing block receipts", "chain", chain.Chain+" "+chain.Name, "block", i, "latest", latest)
		}
		receipts, err := client.BlockReceipts(ctx, rpc.BlockNumberOrHashWithNumber(rpc.BlockNumber(i)))
		if err != nil {
			sugar.Errorw("Failed to fetch latest block receipts", "err", err, "height", i, "chain", chain)
			return err
		}
		if err = FilterEvents(ctx, sugar, q, client, chain, rdb, receipts); err != nil {
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

func FilterEvents(ctx context.Context, sugar *zap.SugaredLogger, q *db.Queries, client *ethclient.Client,
	chain *db.Chain, rdb *redis.Client, receipts utypes.Receipts) error {

	for _, r := range receipts {
		for _, l := range r.Logs { //sugar.Debugw("FilterEvents", "txHash", r.TxHash.Hex(), "l.Address.Hex", l.Address.Hex(), "info", contractType[chain.ID][l.Address.Hex()])

			switch contractType[chain.ID][l.Address.Hex()].Type {
			case db.AssetTypeERC20:
				if err := handleErc20Transfer(ctx, sugar, q, client, chain, rdb, l, false); err != nil {
					sugar.Errorw("handleErc20Transfer", "err", err)
					return err
				}
			case db.AssetTypeERC721:
				if err := handleErc721Transfer(ctx, sugar, q, client, chain, rdb, l, false); err != nil {
					sugar.Errorw("handleErc721Transfer", "err", err)
					return err
				}
			case db.AssetTypeERC1155:
				if l.Topics[0].Hex() == utils.TransferSingleSig {
					if err := handleErc1155TransferSingle(ctx, sugar, q, client, chain, rdb, l, false); err != nil {
						sugar.Errorw("handleErc1155TransferSingle", "err", err)
						return err
					}
				}
				if l.Topics[0].Hex() == utils.TransferBatchSig {
					if err := handleErc1155TransferBatch(ctx, sugar, q, client, chain, rdb, l, false); err != nil {
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

func handleErc20Transfer(ctx context.Context, sugar *zap.SugaredLogger, q *db.Queries, client *ethclient.Client,
	chain *db.Chain, rc *redis.Client, l *utypes.Log, backfill bool) error {

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
	amount := event.Value.String()

	history, err := q.AddOnchainTransaction(ctx, db.AddOnchainTransactionParams{
		From:      event.From.Hex(),
		To:        event.To.Hex(),
		AssetID:   contractType[chain.ID][l.Address.Hex()].ID,
		TokenID:   "0",
		Amount:    amount,
		TxHash:    l.TxHash.Hex(),
		Timestamp: time.Now(),
	})

	if err != nil {
		return err
	}

	if !backfill {
		if err = rdb.SetHistoryCache(ctx, rc, history); err != nil {
			return err
		}
	}

	// Update sender's balances
	balance, err := getErc20BalanceOf(ctx, sugar, client, &l.Address, &event.From)

	if err != nil {
		sugar.Errorw("Failed to get ERC20 balance", "err", err)
	}
	if balance.Cmp(big.NewInt(0)) != 0 {
		if err = q.Add20Asset(ctx, db.Add20AssetParams{
			AssetID: contractType[chain.ID][l.Address.Hex()].ID,
			ChainID: chain.ID,
			Owner:   event.From.Hex(),
			Balance: balance.String(),
		}); err != nil {
			return err
		}
	}

	// Update receiver's balances
	balance, err = getErc20BalanceOf(ctx, sugar, client, &l.Address, &event.To)

	if err != nil {
		sugar.Errorw("Failed to get ERC20 balance", "err", err)
	}
	if balance.Cmp(big.NewInt(0)) != 0 {
		if err = q.Add20Asset(ctx, db.Add20AssetParams{
			AssetID: contractType[chain.ID][l.Address.Hex()].ID,
			ChainID: chain.ID,
			Owner:   event.To.Hex(),
			Balance: balance.String(),
		}); err != nil {
			return err
		}
	}

	return nil
}

func handleErc721Transfer(ctx context.Context, sugar *zap.SugaredLogger, q *db.Queries, client *ethclient.Client,
	chain *db.Chain, rc *redis.Client, l *utypes.Log, backfill bool) error {

	if l.Topics[0].Hex() != utils.TransferEventSig {
		return nil
	}
	// Decode the indexed fields manually
	event := utils.Erc721TransferEvent{
		From:    common.BytesToAddress(l.Topics[1].Bytes()),
		To:      common.BytesToAddress(l.Topics[2].Bytes()),
		TokenID: l.Topics[3].Big(),
	}
	history, err := q.AddOnchainTransaction(ctx, db.AddOnchainTransactionParams{
		From:      event.From.Hex(),
		To:        event.To.Hex(),
		AssetID:   contractType[chain.ID][l.Address.Hex()].ID,
		TokenID:   event.TokenID.String(),
		Amount:    "0",
		TxHash:    l.TxHash.Hex(),
		Timestamp: time.Now(),
	})
	if err != nil {
		return err
	}

	// Cache the new onchain transaction
	if !backfill {
		if err = rdb.SetHistoryCache(ctx, rc, history); err != nil {
			return err
		}
	}

	// Update NFT holder
	uri, err := getErc721TokenURI(ctx, sugar, client, &l.Address, event.TokenID)

	// Get owner of the token
	owner, err := getErc721OwnerOf(ctx, sugar, client, &l.Address, event.TokenID)
	if err != nil {
		sugar.Errorw("Failed to get ERC721 owner", "err", err)
	}

	if err = q.Add721Asset(ctx, db.Add721AssetParams{
		AssetID: contractType[chain.ID][l.Address.Hex()].ID,
		ChainID: chain.ID,
		TokenID: event.TokenID.String(),
		Owner:   owner.Hex(),
		Attributes: sql.NullString{
			String: uri,
			Valid:  true,
		},
	}); err != nil {
		return err
	}

	return nil
}

func handleErc1155TransferBatch(ctx context.Context, sugar *zap.SugaredLogger, q *db.Queries, client *ethclient.Client,
	chain *db.Chain, rc *redis.Client, l *utypes.Log, backfill bool) error {
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
		amount := event.Values[i].String()
		history, err := q.AddOnchainTransaction(ctx, db.AddOnchainTransactionParams{
			From:      event.From.Hex(),
			To:        event.To.Hex(),
			AssetID:   contractType[chain.ID][l.Address.Hex()].ID,
			TokenID:   event.Ids[i].String(),
			Amount:    amount,
			TxHash:    l.TxHash.Hex(),
			Timestamp: time.Now(),
		})
		if err != nil {
			return err
		}

		if !backfill {
			if err = rdb.SetHistoryCache(ctx, rc, history); err != nil {
				return err
			}
		}

		uri, err := getErc1155TokenURI(ctx, sugar, client, &l.Address, event.Ids[i])
		if err != nil {
			sugar.Errorw("Failed to get ERC1155 token URI", "err", err, "tokenID", event.Ids[i])
			return err
		}

		// Update sender's balance
		balance, err := getErc1155BalanceOf(ctx, sugar, client, &l.Address, &event.From, event.Ids[i])
		if err != nil {
			sugar.Errorw("Failed to get ERC1155 balance", "err", err, "tokenID", event.Ids[i])
			return err
		}

		if err = q.Add1155Asset(ctx, db.Add1155AssetParams{
			AssetID: contractType[chain.ID][l.Address.Hex()].ID,
			ChainID: chain.ID,
			TokenID: event.Ids[i].String(),
			Owner:   event.From.Hex(),
			Balance: balance.String(),
			Attributes: sql.NullString{
				String: uri,
				Valid:  true,
			},
		}); err != nil {
			return err
		}

		// Update receiver's balance
		balance, err = getErc1155BalanceOf(ctx, sugar, client, &l.Address, &event.To, event.Ids[i])
		if err != nil {
			sugar.Errorw("Failed to get ERC1155 balance", "err", err, "tokenID", event.Ids[i])
			return err
		}

		if err = q.Add1155Asset(ctx, db.Add1155AssetParams{
			AssetID: contractType[chain.ID][l.Address.Hex()].ID,
			ChainID: chain.ID,
			TokenID: event.Ids[i].String(),
			Owner:   event.To.Hex(),
			Balance: balance.String(),
			Attributes: sql.NullString{
				String: uri,
				Valid:  true,
			},
		}); err != nil {
			return err
		}

	}

	return nil
}

func handleErc1155TransferSingle(ctx context.Context, sugar *zap.SugaredLogger, q *db.Queries, client *ethclient.Client,
	chain *db.Chain, rc *redis.Client, l *utypes.Log, backfill bool) error {

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

	amount := event.Value.String()
	history, err := q.AddOnchainTransaction(ctx, db.AddOnchainTransactionParams{
		From:      event.From.Hex(),
		To:        event.To.Hex(),
		AssetID:   contractType[chain.ID][l.Address.Hex()].ID,
		TokenID:   event.Id.String(),
		Amount:    amount,
		TxHash:    l.TxHash.Hex(),
		Timestamp: time.Now(),
	})
	if err != nil {
		return err
	}

	if !backfill {
		if err = rdb.SetHistoryCache(ctx, rc, history); err != nil {
			return err
		}
	}

	uri, err := getErc1155TokenURI(ctx, sugar, client, &l.Address, event.Id)
	if err != nil {
		sugar.Errorw("Failed to get ERC1155 token URI", "err", err, "tokenID", event.Id)
		return err
	}

	// Update Sender's balance
	balance, err := getErc1155BalanceOf(ctx, sugar, client, &l.Address, &event.From, event.Id)
	if err != nil {
		sugar.Errorw("Failed to get ERC1155 balance", "err", err, "tokenID", event.Id)
		return err
	}

	if err = q.Add1155Asset(ctx, db.Add1155AssetParams{
		AssetID: contractType[chain.ID][l.Address.Hex()].ID,
		ChainID: chain.ID,
		TokenID: event.Id.String(),
		Owner:   event.From.Hex(),
		Balance: balance.String(),
		Attributes: sql.NullString{
			String: uri,
			Valid:  true,
		},
	}); err != nil {
		return err
	}

	// Update Sender's balance
	balance, err = getErc1155BalanceOf(ctx, sugar, client, &l.Address, &event.To, event.Id)
	if err != nil {
		sugar.Errorw("Failed to get ERC1155 balance", "err", err, "tokenID", event.Id)
		return err
	}

	if err = q.Add1155Asset(ctx, db.Add1155AssetParams{
		AssetID: contractType[chain.ID][l.Address.Hex()].ID,
		ChainID: chain.ID,
		TokenID: event.Id.String(),
		Owner:   event.To.Hex(),
		Balance: balance.String(),
		Attributes: sql.NullString{
			String: uri,
			Valid:  true,
		},
	}); err != nil {
		return err
	}

	return nil
}

func getErc721TokenURI(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client,
	contractAddress *common.Address, tokenId *big.Int) (string, error) {
	// Prepare the function call data
	data, err := utils.ERC721ABI.Pack("tokenURI", tokenId)
	if err != nil {
		sugar.Errorf("Failed to pack data for tokenURI: %v", err)
		return "", err
	}

	// Call the contract
	msg := u2u.CallMsg{
		To:   contractAddress,
		Data: data,
	}

	// Execute the call
	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		sugar.Errorf("Failed to call contract: %v", err)
		return "", err
	}

	// Unpack the result to get the token URI
	var tokenURI string
	err = utils.ERC721ABI.UnpackIntoInterface(&tokenURI, "tokenURI", result)
	if err != nil {
		sugar.Errorf("Failed to unpack tokenURI: %v", err)
		return "", err
	}
	return tokenURI, nil
}

func getErc1155TokenURI(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client,
	contractAddress *common.Address, tokenId *big.Int) (string, error) {
	// Prepare the function call data
	data, err := utils.ERC1155ABI.Pack("uri", tokenId)

	if err != nil {
		sugar.Errorf("Failed to pack data for tokenURI: %v", err)
		return "", err
	}

	// Call the contract
	msg := u2u.CallMsg{
		To:   contractAddress,
		Data: data,
	}

	// Execute the call
	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		sugar.Errorf("Failed to call contract: %v", err)
		return "", err
	}

	// Unpack the result to get the token URI
	var tokenURI string
	err = utils.ERC1155ABI.UnpackIntoInterface(&tokenURI, "uri", result)
	if err != nil {
		sugar.Errorf("Failed to unpack tokenURI: %v", err)
		return "", err
	}
	// Replace {id} in the URI template with the actual token ID in hexadecimal form
	tokenIDHex := fmt.Sprintf("%x", tokenId)
	tokenURI = replaceTokenIDPlaceholder(tokenURI, tokenIDHex)
	return tokenURI, nil
}

// replaceTokenIDPlaceholder replaces the "{id}" placeholder with the actual token ID in hexadecimal
func replaceTokenIDPlaceholder(uriTemplate, tokenIDHex string) string {
	return strings.ReplaceAll(uriTemplate, "{id}", tokenIDHex)
}

func retrieveNftMetadata(tokenURI string) ([]byte, error) {
	res, err := http.Get(tokenURI)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(res.Body)
}

func getErc20BalanceOf(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client,
	contractAddress *common.Address, ownerAddress *common.Address) (*big.Int, error) {

	// Prepare the function call data
	data, err := utils.ERC20ABI.Pack("balanceOf", ownerAddress)
	if err != nil {
		sugar.Errorf("Failed to pack data for balanceOf: %v", err)
		return nil, err
	}

	// Call the contract
	msg := u2u.CallMsg{
		To:   contractAddress,
		Data: data,
	}

	// Execute the call
	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		sugar.Errorf("Failed to call contract: %v", err)
		return nil, err
	}

	// Unpack the result to get the balance
	var balance *big.Int
	err = utils.ERC20ABI.UnpackIntoInterface(&balance, "balanceOf", result)
	if err != nil {
		sugar.Errorf("Failed to unpack balanceOf: %v", err)
		return nil, err
	}

	return balance, nil
}

func getErc721OwnerOf(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client,
	contractAddress *common.Address, tokenId *big.Int) (common.Address, error) {

	// Prepare the function call data
	data, err := utils.ERC721ABI.Pack("ownerOf", tokenId)
	if err != nil {
		sugar.Errorf("Failed to pack data for balanceOf: %v", err)
		return common.Address{}, err
	}

	// Call the contract
	msg := u2u.CallMsg{
		To:   contractAddress,
		Data: data,
	}

	// Execute the call
	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		sugar.Errorf("Failed to call contract: %v", err)
		return common.Address{}, err
	}

	// Unpack the result to get the balance
	var owner common.Address
	err = utils.ERC721ABI.UnpackIntoInterface(&owner, "ownerOf", result)

	if err != nil {
		sugar.Errorf("Failed to unpack ownerOf: %v", err)
		return common.Address{}, err
	}

	return owner, nil
}

func getErc1155BalanceOf(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client,
	contractAddress *common.Address, ownerAddress *common.Address, tokenId *big.Int) (*big.Int, error) {
	// Prepare the function call data
	data, err := utils.ERC1155ABI.Pack("balanceOf", ownerAddress, tokenId)
	if err != nil {
		sugar.Errorf("Failed to pack data for balanceOf: %v", err)
		return nil, err
	}

	// Call the contract
	msg := u2u.CallMsg{
		To:   contractAddress,
		Data: data,
	}

	// Execute the call
	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		sugar.Errorf("Failed to call contract: %v", err)
		return nil, err
	}

	// Unpack the result to get the balance
	var balance *big.Int
	err = utils.ERC1155ABI.UnpackIntoInterface(&balance, "balanceOf", result)
	if err != nil {
		sugar.Errorf("Failed to unpack balanceOf: %v", err)
		return nil, err
	}

	return balance, nil
}

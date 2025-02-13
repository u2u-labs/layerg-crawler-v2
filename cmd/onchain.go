package cmd

import (
	"context"
	"math/big"
	"time"

	"github.com/unicornultrafoundation/go-u2u/ethclient"
	"go.uber.org/zap"

	db "github.com/u2u-labs/layerg-crawler/db/sqlc"
	"github.com/u2u-labs/layerg-crawler/generated/router"
)

func StartChainCrawler(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client, q *db.Queries, chain *db.Chain, registry *router.EventRouter) {
	sugar.Infow("Start chain crawler", "chain", chain)
	timer := time.NewTimer(time.Duration(chain.BlockTime) * time.Millisecond)
	defer timer.Stop()
	for range timer.C {
		// Process new blocks
		ProcessLatestBlocks(ctx, sugar, client, q, chain, registry)
		timer.Reset(time.Duration(chain.BlockTime) * time.Millisecond)
	}
}

// func AddBackfillCrawlerTask(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client, q *db.Queries, chain *db.Chain, bf *db.GetCrawlingBackfillCrawlerRow, queueClient *asynq.Client) {
// 	blockRangeScan := int64(config.BackfillBlockRangeScan) * 100
// 	if bf.CurrentBlock%blockRangeScan == 0 {
// 		sugar.Infow("Backfill crawler", "chain", chain, "block", bf.CurrentBlock)
// 	}

// 	timer := time.NewTimer(time.Duration(chain.BlockTime) * time.Millisecond)
// 	defer timer.Stop()

// 	for {
// 		select {
// 		case <-timer.C:
// 			task, err := NewBackfillCollectionTask(bf)
// 			if err != nil {
// 				log.Fatalf("could not create task: %v", err)
// 			}
// 			_, err = queueClient.Enqueue(task)
// 			if err != nil {
// 				log.Fatalf("could not enqueue task: %v", err)
// 			}
// 		}
// 	}

// }

func ProcessLatestBlocks(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client, q *db.Queries, chain *db.Chain, registry *router.EventRouter) error {
	latest, err := client.BlockNumber(ctx)
	if err != nil {
		sugar.Errorw("Failed to fetch latest blocks", "err", err, "chain", chain)
		return err
	}
	sugar.Infow("Processing latest blocks", "chain", chain.Chain+" "+chain.Name, "latest", latest)
	// Process each block between
	for i := chain.LatestBlock + 1; i <= int64(latest); i++ {
		if i%50 == 0 {
			sugar.Infow("Processing block", "chain", chain.Chain+" "+chain.Name, "block", i, "latest", latest)
		}

		block, err := client.BlockByNumber(ctx, big.NewInt(i))
		if err != nil {
			sugar.Errorw("Failed to fetch block", "err", err, "height", i, "chain", chain)
			return err
		}

		// Process each transaction in the block
		for _, tx := range block.Transactions() {
			receipt, err := client.TransactionReceipt(ctx, tx.Hash())
			if err != nil {
				sugar.Errorw("Failed to get receipt", "err", err, "tx", tx.Hash())
				continue
			}

			// Route each log to its handler
			for _, log := range receipt.Logs {
				if err := registry.Route(ctx, log); err != nil {
					sugar.Debugw("Failed to handle event", "err", err, "tx", tx.Hash())
					continue
				}
			}
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

// func FilterEvents(ctx context.Context, sugar *zap.SugaredLogger, q *db.Queries, client *ethclient.Client,
// 	chain *db.Chain, rdb *redis.Client, receipts utypes.Receipts) error {

// 	for _, r := range receipts {
// 		for _, l := range r.Logs {
// 			handlers, err := q.GetContractEventHandlers(ctx, db.GetContractEventHandlersParams{
// 				Address: l.Address.Hex(),
// 				ChainID: chain.ID,
// 			})

// 			if err == nil && len(handlers) > 0 {
// 				registry, err := getContractRegistry(ctx, q, l.Address.Hex(), chain.ID)
// 				if err != nil {
// 					sugar.Errorw("Failed to get contract registry", "err", err)
// 					continue
// 				}

// 				// Get contract ABI to calculate event signatures
// 				contract, err := q.GetContractByAddressAndChain(ctx, db.GetContractByAddressAndChainParams{
// 					Address: l.Address.Hex(),
// 					ChainID: chain.ID,
// 				})
// 				if err != nil {
// 					sugar.Errorw("Failed to get contract", "err", err)
// 					continue
// 				}

// 				contractAbi, err := abi.JSON(strings.NewReader(string(contract.Abi)))
// 				if err != nil {
// 					sugar.Errorw("Failed to parse ABI", "err", err)
// 					continue
// 				}
// 				for _, h := range handlers {
// 					// Get the event definition from ABI
// 					event, exists := contractAbi.Events[h.EventName]
// 					if !exists {
// 						sugar.Errorw("Event not found in ABI", "eventName", h.EventName)
// 						continue
// 					}
// 					sugar.Infow("Processing handlers",
// 						"event", event.ID,
// 						"topicHash", l.Topics[0],
// 					)

// 					// Compare the event signature hash
// 					if l.Topics[0] == event.ID {
// 						handler, ok := registry.Get(h.HandlerName)
// 						if !ok {
// 							sugar.Infow("Handler not found", "handler", h.HandlerName)
// 							continue
// 						}

// 						// Decode and handle event
// 						eventData, err := decodeEvent(contractAbi, l, h.EventName)
// 						if err != nil {
// 							sugar.Errorw("Failed to decode event", "err", err)
// 							continue
// 						}

// 						if err := handler.Handle(ctx, eventData); err != nil {
// 							sugar.Errorw("Failed to handle event", "err", err)
// 						}
// 					}
// 				}
// 			}

// 			// Continue with existing ERC20/721/1155 handling
// 		}
// 	}
// 	return nil
// }

// func handleErc20Transfer(ctx context.Context, sugar *zap.SugaredLogger, q *db.Queries, client *ethclient.Client,
// 	chain *db.Chain, rc *redis.Client, l *utypes.Log) error {

// 	if l.Topics[0].Hex() != utils.TransferEventSig {
// 		return nil
// 	}
// 	// Unpack the log data
// 	var event utils.Erc20TransferEvent

// 	err := utils.ERC20ABI.UnpackIntoInterface(&event, "Transfer", l.Data)
// 	if err != nil {
// 		sugar.Fatalf("Failed to unpack log: %v", err)
// 		return err
// 	}

// 	// Decode the indexed fields manually
// 	event.From = common.BytesToAddress(l.Topics[1].Bytes())
// 	event.To = common.BytesToAddress(l.Topics[2].Bytes())
// 	amount := event.Value.String()

// 	_, err = q.AddOnchainTransaction(ctx, db.AddOnchainTransactionParams{
// 		From:      event.From.Hex(),
// 		To:        event.To.Hex(),
// 		AssetID:   contractType[chain.ID][l.Address.Hex()].ID,
// 		TokenID:   "0",
// 		Amount:    amount,
// 		TxHash:    l.TxHash.Hex(),
// 		Timestamp: time.Now(),
// 	})

// 	if err != nil {
// 		return err
// 	}

// 	// Update sender's balances
// 	balance, err := getErc20BalanceOf(ctx, sugar, client, &l.Address, &event.From)

// 	if err != nil {
// 		sugar.Errorw("Failed to get ERC20 balance", "err", err)
// 	}
// 	if err = q.Add20Asset(ctx, db.Add20AssetParams{
// 		AssetID: contractType[chain.ID][l.Address.Hex()].ID,
// 		ChainID: chain.ID,
// 		Owner:   event.From.Hex(),
// 		Balance: balance.String(),
// 	}); err != nil {
// 		return err
// 	}

// 	// Update receiver's balances
// 	balance, err = getErc20BalanceOf(ctx, sugar, client, &l.Address, &event.To)

// 	if err != nil {
// 		sugar.Errorw("Failed to get ERC20 balance", "err", err)
// 	}
// 	if err = q.Add20Asset(ctx, db.Add20AssetParams{
// 		AssetID: contractType[chain.ID][l.Address.Hex()].ID,
// 		ChainID: chain.ID,
// 		Owner:   event.To.Hex(),
// 		Balance: balance.String(),
// 	}); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func handleErc20BackFill(ctx context.Context, sugar *zap.SugaredLogger, q *db.Queries, client *ethclient.Client,
// 	chain *db.Chain, logs []utypes.Log) error {

// 	// Initialize the AddressSet
// 	addressSet := helpers.NewAddressSet()

// 	if len(logs) == 0 {
// 		return nil
// 	}

// 	var contractAddress *common.Address
// 	for _, l := range logs {
// 		contractAddress = &l.Address
// 		var event utils.Erc20TransferEvent

// 		err := utils.ERC20ABI.UnpackIntoInterface(&event, "Transfer", l.Data)
// 		if err != nil {
// 			sugar.Fatalf("Failed to unpack log: %v", err)
// 			return err
// 		}

// 		if l.Topics[0].Hex() != utils.TransferEventSig {
// 			return nil
// 		}

// 		event.From = common.BytesToAddress(l.Topics[1].Bytes())
// 		event.To = common.BytesToAddress(l.Topics[2].Bytes())
// 		amount := event.Value.String()

// 		_, err = q.AddOnchainTransaction(ctx, db.AddOnchainTransactionParams{
// 			From:      event.From.Hex(),
// 			To:        event.To.Hex(),
// 			AssetID:   contractType[chain.ID][l.Address.Hex()].ID,
// 			TokenID:   "0",
// 			Amount:    amount,
// 			TxHash:    l.TxHash.Hex(),
// 			Timestamp: time.Now(),
// 		})

// 		// adding sender and receiver to the address set
// 		addressSet.AddAddress(event.From)
// 		addressSet.AddAddress(event.To)
// 	}

// 	rpcClient, _ := helpers.InitNewRPCClient(chain.RpcUrl)

// 	addressList := addressSet.GetAddresses()

// 	results := make([]string, len(addressList))
// 	calls := make([]rpc.BatchElem, len(addressList))

// 	for i, addr := range addressList {
// 		// Pack the data for the balanceOf function
// 		data, err := utils.ERC20ABI.Pack("balanceOf", addr)
// 		if err != nil {
// 			sugar.Errorf("Failed to pack data for balanceOf: %v", err)
// 			return err
// 		}

// 		encodedData := "0x" + common.Bytes2Hex(data)

// 		// Append the BatchElem for the eth_call
// 		calls[i] = rpc.BatchElem{
// 			Method: "eth_call",
// 			Args: []interface{}{
// 				map[string]interface{}{
// 					"to":   contractAddress,
// 					"data": encodedData,
// 				},
// 				"latest",
// 			},
// 			Result: &results[i],
// 		}
// 	}

// 	// Execute batch call
// 	if err := rpcClient.BatchCallContext(ctx, calls); err != nil {
// 		log.Fatalf("Failed to execute batch call: %v", err)
// 	}

// 	// Iterate over the results and update the balances
// 	for i, result := range results {
// 		var balance *big.Int

// 		utils.ERC20ABI.UnpackIntoInterface(&balance, "balanceOf", common.FromHex(result))

// 		if err := q.Add20Asset(ctx, db.Add20AssetParams{
// 			AssetID: contractType[chain.ID][contractAddress.Hex()].ID,
// 			ChainID: chain.ID,
// 			Owner:   addressList[i].Hex(),
// 			Balance: balance.String(),
// 		}); err != nil {
// 			return err
// 		}
// 	}

// 	addressSet.Reset()
// 	return nil
// }

// func handleErc721BackFill(ctx context.Context, sugar *zap.SugaredLogger, q *db.Queries, client *ethclient.Client,
// 	chain *db.Chain, logs []utypes.Log) error {

// 	// Initialize the NewTokenIdSet
// 	tokenIdSet := helpers.NewTokenIdSet()

// 	if len(logs) == 0 {
// 		return nil
// 	}

// 	var contractAddress *common.Address
// 	for _, l := range logs {
// 		contractAddress = &l.Address

// 		// Decode the indexed fields manually
// 		event := utils.Erc721TransferEvent{
// 			From:    common.BytesToAddress(l.Topics[1].Bytes()),
// 			To:      common.BytesToAddress(l.Topics[2].Bytes()),
// 			TokenID: l.Topics[3].Big(),
// 		}
// 		_, err := q.AddOnchainTransaction(ctx, db.AddOnchainTransactionParams{
// 			From:      event.From.Hex(),
// 			To:        event.To.Hex(),
// 			AssetID:   contractType[chain.ID][l.Address.Hex()].ID,
// 			TokenID:   event.TokenID.String(),
// 			Amount:    "0",
// 			TxHash:    l.TxHash.Hex(),
// 			Timestamp: time.Now(),
// 		})
// 		if err != nil {
// 			return err
// 		}

// 		// adding token Id
// 		tokenIdSet.AddTokenId(event.TokenID)
// 	}

// 	rpcClient, _ := helpers.InitNewRPCClient(chain.RpcUrl)

// 	tokenIdList := tokenIdSet.GetTokenIds()

// 	results := make([]string, len(tokenIdList)*2)
// 	calls := make([]rpc.BatchElem, len(tokenIdList)*2)

// 	for i, tokenId := range tokenIdList {
// 		// Pack the data for the tokenURI function
// 		data, err := utils.ERC721ABI.Pack("tokenURI", tokenId)
// 		if err != nil {
// 			sugar.Errorf("Failed to pack data for tokenURI: %v", err)
// 			return err
// 		}

// 		encodedUriData := "0x" + common.Bytes2Hex(data)

// 		// Append the BatchElem for the eth_call
// 		calls[2*i] = rpc.BatchElem{
// 			Method: "eth_call",
// 			Args: []interface{}{
// 				map[string]interface{}{
// 					"to":   contractAddress,
// 					"data": encodedUriData,
// 				},
// 				"latest",
// 			},
// 			Result: &results[2*i],
// 		}

// 		// Pack the data for the ownerOf function
// 		ownerData, err := utils.ERC721ABI.Pack("ownerOf", tokenId)
// 		if err != nil {
// 			sugar.Errorf("Failed to pack data for ownerOf: %v", err)
// 			return err
// 		}

// 		encodedOwnerData := "0x" + common.Bytes2Hex(ownerData)

// 		// Append the BatchElem for the eth_call
// 		calls[2*i+1] = rpc.BatchElem{
// 			Method: "eth_call",
// 			Args: []interface{}{
// 				map[string]interface{}{
// 					"to":   contractAddress,
// 					"data": encodedOwnerData,
// 				},
// 				"latest",
// 			},
// 			Result: &results[2*i+1],
// 		}

// 	}

// 	// Execute batch call
// 	if err := rpcClient.BatchCallContext(ctx, calls); err != nil {
// 		log.Fatalf("Failed to execute batch call: %v", err)
// 	}

// 	// Iterate over the results and update the balances
// 	for i := 0; i < len(results); i += 2 {
// 		var uri string
// 		var owner common.Address
// 		utils.ERC721ABI.UnpackIntoInterface(&uri, "tokenURI", common.FromHex(results[i]))
// 		utils.ERC721ABI.UnpackIntoInterface(&owner, "ownerOf", common.FromHex(results[i+1]))

// 		if err := q.Add721Asset(ctx, db.Add721AssetParams{
// 			AssetID: contractType[chain.ID][contractAddress.Hex()].ID,
// 			ChainID: chain.ID,
// 			TokenID: tokenIdList[i/2].String(),
// 			Owner:   owner.Hex(),
// 			Attributes: sql.NullString{
// 				String: uri,
// 				Valid:  true,
// 			},
// 		}); err != nil {
// 			return err
// 		}
// 	}

// 	tokenIdSet.Reset()
// 	return nil
// }

// func handleErc1155Backfill(ctx context.Context, sugar *zap.SugaredLogger, q *db.Queries, client *ethclient.Client,
// 	chain *db.Chain, logs []utypes.Log) error {

// 	// Initialize the NewTokenIdSet
// 	tokenIdContractAddressSet := helpers.NewTokenIdContractAddressSet()

// 	if len(logs) == 0 {
// 		return nil
// 	}

// 	var contractAddress *common.Address
// 	for _, l := range logs {
// 		contractAddress = &l.Address
// 		if l.Topics[0].Hex() == utils.TransferSingleSig {
// 			// handleTransferSingle

// 			// Decode TransferSingle log
// 			var event utils.Erc1155TransferSingleEvent
// 			err := utils.ERC1155ABI.UnpackIntoInterface(&event, "TransferSingle", l.Data)
// 			if err != nil {
// 				sugar.Errorw("Failed to unpack TransferSingle log:", "err", err)
// 			}

// 			// Decode the indexed fields for TransferSingle
// 			event.Operator = common.BytesToAddress(l.Topics[1].Bytes())
// 			event.From = common.BytesToAddress(l.Topics[2].Bytes())
// 			event.To = common.BytesToAddress(l.Topics[3].Bytes())

// 			amount := event.Value.String()
// 			_, err = q.AddOnchainTransaction(ctx, db.AddOnchainTransactionParams{
// 				From:      event.From.Hex(),
// 				To:        event.To.Hex(),
// 				AssetID:   contractType[chain.ID][l.Address.Hex()].ID,
// 				TokenID:   event.Id.String(),
// 				Amount:    amount,
// 				TxHash:    l.TxHash.Hex(),
// 				Timestamp: time.Now(),
// 			})
// 			if err != nil {
// 				return err
// 			}

// 			// adding data to set
// 			tokenIdContractAddressSet.AddTokenIdContractAddress(event.Id, event.From.Hex())
// 			tokenIdContractAddressSet.AddTokenIdContractAddress(event.Id, event.To.Hex())
// 		}

// 		if l.Topics[0].Hex() == utils.TransferBatchSig {
// 			var event utils.Erc1155TransferBatchEvent
// 			err := utils.ERC1155ABI.UnpackIntoInterface(&event, "TransferBatch", l.Data)
// 			if err != nil {
// 				sugar.Errorw("Failed to unpack TransferBatch log:", "err", err)
// 			}

// 			// Decode the indexed fields for TransferBatch
// 			event.Operator = common.BytesToAddress(l.Topics[1].Bytes())
// 			event.From = common.BytesToAddress(l.Topics[2].Bytes())
// 			event.To = common.BytesToAddress(l.Topics[3].Bytes())

// 			for i := range event.Ids {
// 				amount := event.Values[i].String()
// 				_, err := q.AddOnchainTransaction(ctx, db.AddOnchainTransactionParams{
// 					From:      event.From.Hex(),
// 					To:        event.To.Hex(),
// 					AssetID:   contractType[chain.ID][l.Address.Hex()].ID,
// 					TokenID:   event.Ids[i].String(),
// 					Amount:    amount,
// 					TxHash:    l.TxHash.Hex(),
// 					Timestamp: time.Now(),
// 				})
// 				if err != nil {
// 					return err
// 				}

// 				// adding data to set
// 				tokenIdContractAddressSet.AddTokenIdContractAddress(event.Ids[i], event.From.Hex())
// 				tokenIdContractAddressSet.AddTokenIdContractAddress(event.Ids[i], event.To.Hex())
// 			}
// 		}
// 	}

// 	rpcClient, _ := helpers.InitNewRPCClient(chain.RpcUrl)

// 	tokenIdList := tokenIdContractAddressSet.GetTokenIdContractAddressses()

// 	results := make([]string, len(tokenIdList)*2)
// 	calls := make([]rpc.BatchElem, len(tokenIdList)*2)

// 	for i, pairData := range tokenIdList {
// 		tokenId := pairData.TokenId
// 		ownerAddress := common.HexToAddress(pairData.ContractAddress)

// 		// Pack the data for the tokenURI function
// 		data, err := utils.ERC1155ABI.Pack("uri", tokenId)
// 		if err != nil {
// 			sugar.Errorf("Failed to pack data for tokenURI: %v", err)
// 			return err
// 		}

// 		encodedUriData := "0x" + common.Bytes2Hex(data)

// 		// Append the BatchElem for the eth_call
// 		calls[2*i] = rpc.BatchElem{
// 			Method: "eth_call",
// 			Args: []interface{}{
// 				map[string]interface{}{
// 					"to":   contractAddress,
// 					"data": encodedUriData,
// 				},
// 				"latest",
// 			},
// 			Result: &results[2*i],
// 		}

// 		// 	// Pack the data for the ownerOf function
// 		ownerData, err := utils.ERC1155ABI.Pack("balanceOf", ownerAddress, tokenId)
// 		if err != nil {
// 			sugar.Errorf("Failed to pack data for balanceOf: %v", err)
// 			return err
// 		}

// 		encodedBalanceData := "0x" + common.Bytes2Hex(ownerData)

// 		// 	// Append the BatchElem for the eth_call
// 		calls[2*i+1] = rpc.BatchElem{
// 			Method: "eth_call",
// 			Args: []interface{}{
// 				map[string]interface{}{
// 					"to":   contractAddress,
// 					"data": encodedBalanceData,
// 				},
// 				"latest",
// 			},
// 			Result: &results[2*i+1],
// 		}
// 	}

// 	// // Execute batch call
// 	if err := rpcClient.BatchCallContext(ctx, calls); err != nil {
// 		log.Fatalf("Failed to execute batch call: %v", err)
// 	}

// 	// // Iterate over the results and update the balances
// 	for i := 0; i < len(results); i += 2 {
// 		var uri string
// 		var balance *big.Int
// 		utils.ERC1155ABI.UnpackIntoInterface(&uri, "uri", common.FromHex(results[i]))
// 		utils.ERC1155ABI.UnpackIntoInterface(&balance, "balanceOf", common.FromHex(results[i+1]))

// 		if err := q.Add1155Asset(ctx, db.Add1155AssetParams{
// 			AssetID: contractType[chain.ID][contractAddress.Hex()].ID,
// 			ChainID: chain.ID,
// 			TokenID: tokenIdList[i/2].TokenId.String(),
// 			Owner:   tokenIdList[i/2].ContractAddress,
// 			Attributes: sql.NullString{
// 				String: uri,
// 				Valid:  true,
// 			},
// 		}); err != nil {
// 			return err
// 		}
// 	}

// 	tokenIdContractAddressSet.Reset()
// 	return nil
// }

// func handleErc721Transfer(ctx context.Context, sugar *zap.SugaredLogger, q *db.Queries, client *ethclient.Client,
// 	chain *db.Chain, rc *redis.Client, l *utypes.Log) error {

// 	if l.Topics[0].Hex() != utils.TransferEventSig {
// 		return nil
// 	}
// 	// Decode the indexed fields manually
// 	event := utils.Erc721TransferEvent{
// 		From:    common.BytesToAddress(l.Topics[1].Bytes()),
// 		To:      common.BytesToAddress(l.Topics[2].Bytes()),
// 		TokenID: l.Topics[3].Big(),
// 	}
// 	_, err := q.AddOnchainTransaction(ctx, db.AddOnchainTransactionParams{
// 		From:      event.From.Hex(),
// 		To:        event.To.Hex(),
// 		AssetID:   contractType[chain.ID][l.Address.Hex()].ID,
// 		TokenID:   event.TokenID.String(),
// 		Amount:    "0",
// 		TxHash:    l.TxHash.Hex(),
// 		Timestamp: time.Now(),
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	// Update NFT holder
// 	uri, err := getErc721TokenURI(ctx, sugar, client, &l.Address, event.TokenID)

// 	// Get owner of the token
// 	owner, err := getErc721OwnerOf(ctx, sugar, client, &l.Address, event.TokenID)
// 	if err != nil {

// 		sugar.Errorw("Failed to get ERC721 owner", "err", err, "tokenID", event.TokenID, "contract", l.Address.Hex())
// 	}

// 	if err != nil {

// 		sugar.Errorw("Failed to get ERC721 owner", "err", err, "tokenID", event.TokenID, "contract", l.Address.Hex())
// 	}
// 	//

// 	if err = q.Add721Asset(ctx, db.Add721AssetParams{
// 		AssetID: contractType[chain.ID][l.Address.Hex()].ID,
// 		ChainID: chain.ID,
// 		TokenID: event.TokenID.String(),
// 		Owner:   owner.Hex(),
// 		Attributes: sql.NullString{
// 			String: uri,
// 			Valid:  true,
// 		},
// 	}); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func addingErc721Elem(ctx context.Context, sugar *zap.SugaredLogger, q *db.Queries, client *ethclient.Client,
// 	chain *db.Chain, rc *redis.Client, l *utypes.Log, backfill bool) error {

// 	if l.Topics[0].Hex() != utils.TransferEventSig {
// 		return nil
// 	}
// 	// Decode the indexed fields manually
// 	event := utils.Erc721TransferEvent{
// 		From:    common.BytesToAddress(l.Topics[1].Bytes()),
// 		To:      common.BytesToAddress(l.Topics[2].Bytes()),
// 		TokenID: l.Topics[3].Big(),
// 	}
// 	history, err := q.AddOnchainTransaction(ctx, db.AddOnchainTransactionParams{
// 		From:      event.From.Hex(),
// 		To:        event.To.Hex(),
// 		AssetID:   contractType[chain.ID][l.Address.Hex()].ID,
// 		TokenID:   event.TokenID.String(),
// 		Amount:    "0",
// 		TxHash:    l.TxHash.Hex(),
// 		Timestamp: time.Now(),
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	// Cache the new onchain transaction
// 	if !backfill {
// 		if err = rdb.SetHistoryCache(ctx, rc, history); err != nil {
// 			return err
// 		}
// 	}

// 	// Update NFT holder
// 	uri, err := getErc721TokenURI(ctx, sugar, client, &l.Address, event.TokenID)

// 	// Get owner of the token
// 	owner, err := getErc721OwnerOf(ctx, sugar, client, &l.Address, event.TokenID)
// 	if err != nil {

// 		sugar.Errorw("Failed to get ERC721 owner", "err", err, "tokenID", event.TokenID, "contract", l.Address.Hex())
// 	}

// 	if err != nil {

// 		sugar.Errorw("Failed to get ERC721 owner", "err", err, "tokenID", event.TokenID, "contract", l.Address.Hex())
// 	}
// 	//

// 	if err = q.Add721Asset(ctx, db.Add721AssetParams{
// 		AssetID: contractType[chain.ID][l.Address.Hex()].ID,
// 		ChainID: chain.ID,
// 		TokenID: event.TokenID.String(),
// 		Owner:   owner.Hex(),
// 		Attributes: sql.NullString{
// 			String: uri,
// 			Valid:  true,
// 		},
// 	}); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func handleErc1155TransferBatch(ctx context.Context, sugar *zap.SugaredLogger, q *db.Queries, client *ethclient.Client,
// 	chain *db.Chain, rc *redis.Client, l *utypes.Log) error {
// 	// Decode TransferBatch log
// 	var event utils.Erc1155TransferBatchEvent
// 	err := utils.ERC1155ABI.UnpackIntoInterface(&event, "TransferBatch", l.Data)
// 	if err != nil {
// 		sugar.Errorw("Failed to unpack TransferBatch log:", "err", err)
// 	}

// 	// Decode the indexed fields for TransferBatch
// 	event.Operator = common.BytesToAddress(l.Topics[1].Bytes())
// 	event.From = common.BytesToAddress(l.Topics[2].Bytes())
// 	event.To = common.BytesToAddress(l.Topics[3].Bytes())

// 	for i := range event.Ids {
// 		amount := event.Values[i].String()
// 		_, err := q.AddOnchainTransaction(ctx, db.AddOnchainTransactionParams{
// 			From:      event.From.Hex(),
// 			To:        event.To.Hex(),
// 			AssetID:   contractType[chain.ID][l.Address.Hex()].ID,
// 			TokenID:   event.Ids[i].String(),
// 			Amount:    amount,
// 			TxHash:    l.TxHash.Hex(),
// 			Timestamp: time.Now(),
// 		})
// 		if err != nil {
// 			return err
// 		}

// 		uri, err := getErc1155TokenURI(ctx, sugar, client, &l.Address, event.Ids[i])
// 		if err != nil {
// 			sugar.Errorw("Failed to get ERC1155 token URI", "err", err, "tokenID", event.Ids[i])
// 			return err
// 		}

// 		// Update sender's balance
// 		balance, err := getErc1155BalanceOf(ctx, sugar, client, &l.Address, &event.From, event.Ids[i])
// 		if err != nil {
// 			sugar.Errorw("Failed to get ERC1155 balance", "err", err, "tokenID", event.Ids[i])
// 			return err
// 		}

// 		if err = q.Add1155Asset(ctx, db.Add1155AssetParams{
// 			AssetID: contractType[chain.ID][l.Address.Hex()].ID,
// 			ChainID: chain.ID,
// 			TokenID: event.Ids[i].String(),
// 			Owner:   event.From.Hex(),
// 			Balance: balance.String(),
// 			Attributes: sql.NullString{
// 				String: uri,
// 				Valid:  true,
// 			},
// 		}); err != nil {
// 			return err
// 		}

// 		// Update receiver's balance
// 		balance, err = getErc1155BalanceOf(ctx, sugar, client, &l.Address, &event.To, event.Ids[i])
// 		if err != nil {
// 			sugar.Errorw("Failed to get ERC1155 balance", "err", err, "tokenID", event.Ids[i])
// 			return err
// 		}

// 		if err = q.Add1155Asset(ctx, db.Add1155AssetParams{
// 			AssetID: contractType[chain.ID][l.Address.Hex()].ID,
// 			ChainID: chain.ID,
// 			TokenID: event.Ids[i].String(),
// 			Owner:   event.To.Hex(),
// 			Balance: balance.String(),
// 			Attributes: sql.NullString{
// 				String: uri,
// 				Valid:  true,
// 			},
// 		}); err != nil {
// 			return err
// 		}

// 	}

// 	return nil
// }

// func handleErc1155TransferSingle(ctx context.Context, sugar *zap.SugaredLogger, q *db.Queries, client *ethclient.Client,
// 	chain *db.Chain, rc *redis.Client, l *utypes.Log) error {

// 	// Decode TransferSingle log
// 	var event utils.Erc1155TransferSingleEvent
// 	err := utils.ERC1155ABI.UnpackIntoInterface(&event, "TransferSingle", l.Data)
// 	if err != nil {
// 		sugar.Errorw("Failed to unpack TransferSingle log:", "err", err)
// 	}

// 	// Decode the indexed fields for TransferSingle
// 	event.Operator = common.BytesToAddress(l.Topics[1].Bytes())
// 	event.From = common.BytesToAddress(l.Topics[2].Bytes())
// 	event.To = common.BytesToAddress(l.Topics[3].Bytes())

// 	amount := event.Value.String()
// 	_, err = q.AddOnchainTransaction(ctx, db.AddOnchainTransactionParams{
// 		From:      event.From.Hex(),
// 		To:        event.To.Hex(),
// 		AssetID:   contractType[chain.ID][l.Address.Hex()].ID,
// 		TokenID:   event.Id.String(),
// 		Amount:    amount,
// 		TxHash:    l.TxHash.Hex(),
// 		Timestamp: time.Now(),
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	uri, err := getErc1155TokenURI(ctx, sugar, client, &l.Address, event.Id)
// 	if err != nil {
// 		sugar.Errorw("Failed to get ERC1155 token URI", "err", err, "tokenID", event.Id)
// 		return err
// 	}

// 	// Update Sender's balance
// 	balance, err := getErc1155BalanceOf(ctx, sugar, client, &l.Address, &event.From, event.Id)
// 	if err != nil {
// 		sugar.Errorw("Failed to get ERC1155 balance", "err", err, "tokenID", event.Id)
// 		return err
// 	}

// 	if err = q.Add1155Asset(ctx, db.Add1155AssetParams{
// 		AssetID: contractType[chain.ID][l.Address.Hex()].ID,
// 		ChainID: chain.ID,
// 		TokenID: event.Id.String(),
// 		Owner:   event.From.Hex(),
// 		Balance: balance.String(),
// 		Attributes: sql.NullString{
// 			String: uri,
// 			Valid:  true,
// 		},
// 	}); err != nil {
// 		return err
// 	}

// 	// Update Sender's balance
// 	balance, err = getErc1155BalanceOf(ctx, sugar, client, &l.Address, &event.To, event.Id)
// 	if err != nil {
// 		sugar.Errorw("Failed to get ERC1155 balance", "err", err, "tokenID", event.Id)
// 		return err
// 	}

// 	if err = q.Add1155Asset(ctx, db.Add1155AssetParams{
// 		AssetID: contractType[chain.ID][l.Address.Hex()].ID,
// 		ChainID: chain.ID,
// 		TokenID: event.Id.String(),
// 		Owner:   event.To.Hex(),
// 		Balance: balance.String(),
// 		Attributes: sql.NullString{
// 			String: uri,
// 			Valid:  true,
// 		},
// 	}); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func getErc721TokenURI(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client,
// 	contractAddress *common.Address, tokenId *big.Int) (string, error) {
// 	// Prepare the function call data
// 	data, err := utils.ERC721ABI.Pack("tokenURI", tokenId)
// 	if err != nil {
// 		sugar.Errorf("Failed to pack data for tokenURI: %v", err)
// 		return "", err
// 	}

// 	// Call the contract
// 	msg := u2u.CallMsg{
// 		To:   contractAddress,
// 		Data: data,
// 	}

// 	// Execute the call
// 	result, err := client.CallContract(context.Background(), msg, nil)
// 	if err != nil {
// 		sugar.Errorf("Failed to call contract: %v", err)
// 		return "", err
// 	}

// 	// Unpack the result to get the token URI
// 	var tokenURI string
// 	err = utils.ERC721ABI.UnpackIntoInterface(&tokenURI, "tokenURI", result)
// 	if err != nil {
// 		sugar.Errorf("Failed to unpack tokenURI: %v", err)
// 		return "", err
// 	}
// 	return tokenURI, nil
// }

// func getErc1155TokenURI(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client,
// 	contractAddress *common.Address, tokenId *big.Int) (string, error) {
// 	// Prepare the function call data
// 	data, err := utils.ERC1155ABI.Pack("uri", tokenId)

// 	if err != nil {
// 		sugar.Errorf("Failed to pack data for tokenURI: %v", err)
// 		return "", err
// 	}

// 	// Call the contract
// 	msg := u2u.CallMsg{
// 		To:   contractAddress,
// 		Data: data,
// 	}

// 	// Execute the call
// 	result, err := client.CallContract(context.Background(), msg, nil)
// 	if err != nil {
// 		sugar.Errorf("Failed to call contract: %v", err)
// 		return "", err
// 	}

// 	// Unpack the result to get the token URI
// 	var tokenURI string
// 	err = utils.ERC1155ABI.UnpackIntoInterface(&tokenURI, "uri", result)
// 	if err != nil {
// 		sugar.Errorf("Failed to unpack tokenURI: %v", err)
// 		return "", err
// 	}
// 	// Replace {id} in the URI template with the actual token ID in hexadecimal form
// 	tokenIDHex := fmt.Sprintf("%x", tokenId)
// 	tokenURI = replaceTokenIDPlaceholder(tokenURI, tokenIDHex)
// 	return tokenURI, nil
// }

// // replaceTokenIDPlaceholder replaces the "{id}" placeholder with the actual token ID in hexadecimal
// func replaceTokenIDPlaceholder(uriTemplate, tokenIDHex string) string {
// 	return strings.ReplaceAll(uriTemplate, "{id}", tokenIDHex)
// }

// func retrieveNftMetadata(tokenURI string) ([]byte, error) {
// 	res, err := http.Get(tokenURI)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return io.ReadAll(res.Body)
// }

// func getErc20BalanceOf(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client,
// 	contractAddress *common.Address, ownerAddress *common.Address) (*big.Int, error) {

// 	// Prepare the function call data
// 	data, err := utils.ERC20ABI.Pack("balanceOf", ownerAddress)
// 	if err != nil {
// 		sugar.Errorf("Failed to pack data for balanceOf: %v", err)
// 		return nil, err
// 	}

// 	// Call the contract
// 	msg := u2u.CallMsg{
// 		To:   contractAddress,
// 		Data: data,
// 	}

// 	// Execute the call
// 	result, err := client.CallContract(context.Background(), msg, nil)
// 	if err != nil {
// 		sugar.Errorf("Failed to call contract: %v", err)
// 		return nil, err
// 	}

// 	// Unpack the result to get the balance
// 	var balance *big.Int
// 	err = utils.ERC20ABI.UnpackIntoInterface(&balance, "balanceOf", result)
// 	if err != nil {
// 		sugar.Errorf("Failed to unpack balanceOf: %v", err)
// 		return nil, err
// 	}

// 	return balance, nil
// }

// func getErc721OwnerOf(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client,
// 	contractAddress *common.Address, tokenId *big.Int) (common.Address, error) {

// 	// Prepare the function call data
// 	data, err := utils.ERC721ABI.Pack("ownerOf", tokenId)
// 	if err != nil {
// 		sugar.Errorf("Failed to pack data for balanceOf: %v", err)
// 		return common.Address{}, err
// 	}

// 	// Call the contract
// 	msg := u2u.CallMsg{
// 		To:   contractAddress,
// 		Data: data,
// 	}

// 	// Execute the call
// 	result, err := client.CallContract(context.Background(), msg, nil)
// 	if err != nil {
// 		sugar.Errorf("Failed to call contract: %v", err)
// 		return common.Address{}, err
// 	}

// 	// Unpack the result to get the balance
// 	var owner common.Address
// 	err = utils.ERC721ABI.UnpackIntoInterface(&owner, "ownerOf", result)

// 	if err != nil {
// 		sugar.Errorf("Failed to unpack ownerOf: %v", err)
// 		return common.Address{}, err
// 	}

// 	return owner, nil
// }

// func getErc1155BalanceOf(ctx context.Context, sugar *zap.SugaredLogger, client *ethclient.Client,
// 	contractAddress *common.Address, ownerAddress *common.Address, tokenId *big.Int) (*big.Int, error) {
// 	// Prepare the function call data
// 	data, err := utils.ERC1155ABI.Pack("balanceOf", ownerAddress, tokenId)
// 	if err != nil {
// 		sugar.Errorf("Failed to pack data for balanceOf: %v", err)
// 		return nil, err
// 	}

// 	// Call the contract
// 	msg := u2u.CallMsg{
// 		To:   contractAddress,
// 		Data: data,
// 	}

// 	// Execute the call
// 	result, err := client.CallContract(context.Background(), msg, nil)
// 	if err != nil {
// 		sugar.Errorf("Failed to call contract: %v", err)
// 		return nil, err
// 	}

// 	// Unpack the result to get the balance
// 	var balance *big.Int
// 	err = utils.ERC1155ABI.UnpackIntoInterface(&balance, "balanceOf", result)
// 	if err != nil {
// 		sugar.Errorf("Failed to unpack balanceOf: %v", err)
// 		return nil, err
// 	}

// 	return balance, nil
// }

package handlers

import (
	"context"
	"fmt"
	"math/big"

	"go.uber.org/zap"

	"github.com/google/uuid"
	graphqldb "github.com/u2u-labs/layerg-crawler/db/graphqldb"
	db "github.com/u2u-labs/layerg-crawler/db/sqlc"
	"github.com/u2u-labs/layerg-crawler/generated/eventhandlers"
)

type TransferHandler struct {
	*BaseHandler
}

type MetadataUpdateHandler struct {
	*BaseHandler
}

type TransferSingleHandler struct {
	*BaseHandler
}

type TransferBatchHandler struct {
	*BaseHandler
}

type URIHandler struct {
	*BaseHandler
}

func NewTransferHandler(queries *db.Queries, gqlQueries *graphqldb.Queries, chainID int32, logger *zap.SugaredLogger) *TransferHandler {
	return &TransferHandler{
		BaseHandler: &BaseHandler{
			Queries: queries,
			GQL:     gqlQueries,
			ChainID: chainID,
			Logger:  logger,
		},
	}
}

func NewMetadataUpdateHandler(queries *db.Queries, gqlQueries *graphqldb.Queries, chainID int32, logger *zap.SugaredLogger) *MetadataUpdateHandler {
	return &MetadataUpdateHandler{
		BaseHandler: &BaseHandler{
			Queries: queries,
			GQL:     gqlQueries,
			ChainID: chainID,
			Logger:  logger,
		},
	}
}

func NewTransferSingleHandler(queries *db.Queries, gqlQueries *graphqldb.Queries, chainID int32, logger *zap.SugaredLogger) *TransferSingleHandler {
	return &TransferSingleHandler{
		BaseHandler: &BaseHandler{
			Queries: queries,
			GQL:     gqlQueries,
			ChainID: chainID,
			Logger:  logger,
		},
	}
}

func NewTransferBatchHandler(queries *db.Queries, gqlQueries *graphqldb.Queries, chainID int32, logger *zap.SugaredLogger) *TransferBatchHandler {
	return &TransferBatchHandler{
		BaseHandler: &BaseHandler{
			Queries: queries,
			GQL:     gqlQueries,
			ChainID: chainID,
			Logger:  logger,
		},
	}
}

func NewURIHandler(queries *db.Queries, gqlQueries *graphqldb.Queries, chainID int32, logger *zap.SugaredLogger) *URIHandler {
	return &URIHandler{
		BaseHandler: &BaseHandler{
			Queries: queries,
			GQL:     gqlQueries,
			ChainID: chainID,
			Logger:  logger,
		},
	}
}

func (h *TransferHandler) HandleTransfer(ctx context.Context, event *eventhandlers.Transfer) error {
	// Handle ERC721 transfer
	fromUser, err := h.GQL.GetOrCreateUser(ctx, event.From.Hex())
	if err != nil {
		h.Logger.Errorw("Failed to ensure 'from' user exists",
			"err", err,
			"address", event.From.Hex(),
		)
		return fmt.Errorf("failed to ensure 'from' user exists: %w", err)
	}
	// h.AddOperation("User", fromUser, event.Raw.BlockHash.Hex(), event.Raw.BlockNumber)

	toUser, err := h.GQL.GetOrCreateUser(ctx, event.To.Hex())
	if err != nil {
		h.Logger.Errorw("Failed to ensure 'to' user exists",
			"err", err,
			"address", event.To.Hex(),
		)
		return fmt.Errorf("failed to ensure 'to' user exists: %w", err)
	}
	// h.AddOperation("User", toUser, event.Raw.BlockHash.Hex(), event.Raw.BlockNumber)

	tokenID := event.TokenId.String()

	// Create or update item
	item, err := h.GQL.GetItemByTokenId(ctx, tokenID)
	if err != nil {
		// Create new item if it doesn't exist
		item, err = h.GQL.CreateItem(ctx, graphqldb.CreateItemParams{
			ID:       uuid.New().String(),
			TokenID:  tokenID,
			TokenUri: "",
			Standard: "ERC721",
		})
		if err != nil {
			return fmt.Errorf("failed to create item: %w", err)
		}
	}

	// Update balances
	if event.From.Hex() != "0x0000000000000000000000000000000000000000" {
		// Remove balance from sender
		_, err = h.GQL.UpsertBalance(ctx, graphqldb.UpsertBalanceParams{
			ID:        uuid.New().String(),
			ItemID:    item.ID,
			OwnerID:   fromUser.ID,
			Value:     "0",
			UpdatedAt: fmt.Sprintf("%d", event.Raw.BlockNumber),
			Contract:  event.Raw.Address.Hex(),
		})
		if err != nil {
			return fmt.Errorf("failed to update sender balance: %w", err)
		}
	}

	// Add balance to receiver
	_, err = h.GQL.UpsertBalance(ctx, graphqldb.UpsertBalanceParams{
		ID:        uuid.New().String(),
		ItemID:    item.ID,
		OwnerID:   toUser.ID,
		Value:     "1",
		UpdatedAt: fmt.Sprintf("%d", event.Raw.BlockNumber),
		Contract:  event.Raw.Address.Hex(),
	})
	if err != nil {
		return fmt.Errorf("failed to update 721 receiver balance: %w", err)
	}

	h.Logger.Infow("Transfer event processed",
		"tokenId", tokenID,
		"from", fromUser.ID,
		"to", toUser.ID,
		"contract", event.Raw.Address.Hex(),
		"tx", event.Raw.TxHash.Hex(),
	)

	if err := h.SubmitToDA(); err != nil {
		h.Logger.Errorw("Failed to submit to DA", "error", err)
	}

	return nil
}

func (h *MetadataUpdateHandler) HandleMetadataUpdate(ctx context.Context, event *eventhandlers.MetadataUpdate) error {
	actor := event.Raw.TxHash.Hex()

	actorUser, err := h.GQL.GetOrCreateUser(ctx, actor)
	if err != nil {
		h.Logger.Errorw("Failed to ensure actor exists",
			"err", err,
			"address", actor,
		)
		return fmt.Errorf("failed to ensure actor exists: %w", err)
	}

	// h.AddOperation("User", actorUser, event.Raw.BlockHash.Hex(), event.Raw.BlockNumber)

	tokenID := event.TokenId.String()
	_, err = h.GQL.CreateMetadataUpdateRecord(ctx, graphqldb.CreateMetadataUpdateRecordParams{
		ID:      uuid.New().String(),
		TokenID: tokenID,
		ActorID: actorUser.ID,
	})

	if err != nil {
		h.Logger.Errorw("Failed to create metadata update record",
			"err", err,
			"tokenId", tokenID,
			"actor", actor,
			"tx", event.Raw.TxHash.Hex(),
		)
		return fmt.Errorf("failed to create metadata update record: %w", err)
	}

	// Update the item's token URI if it exists
	items, err := h.GQL.ListItem(ctx)
	if err != nil {
		return fmt.Errorf("failed to list items: %w", err)
	}

	for _, item := range items {
		if item.TokenID == tokenID {

			newTokenUri := "" // TODO: Implement fetching token URI from contract

			_, err := h.GQL.UpdateItem(ctx, graphqldb.UpdateItemParams{
				ID:       item.ID,
				TokenID:  item.TokenID,
				TokenUri: newTokenUri,
			})
			if err != nil {
				h.Logger.Errorw("Failed to update item token URI",
					"err", err,
					"tokenId", tokenID,
					"newUri", newTokenUri,
				)
				return fmt.Errorf("failed to update item token URI: %w", err)
			}
			// h.AddOperation("Item", itemResp, event.Raw.BlockHash.Hex(), event.Raw.BlockNumber)
			break
		}
	}

	if err := h.SubmitToDA(); err != nil {
		h.Logger.Errorw("Failed to submit to DA", "error", err)
	}

	h.Logger.Infow("Metadata update event processed",
		"tokenId", tokenID,
		"actor", actor,
		"contract", event.Raw.Address.Hex(),
		"tx", event.Raw.TxHash.Hex(),
	)

	return nil
}

func (h *TransferSingleHandler) HandleTransferSingle(ctx context.Context, event *eventhandlers.TransferSingle) error {
	fromUser, err := h.GQL.GetOrCreateUser(ctx, event.From.Hex())
	if err != nil {
		return fmt.Errorf("failed to ensure 'from' user exists: %w", err)
	}
	// h.AddOperation("User", fromUser, event.Raw.BlockHash.Hex(), event.Raw.BlockNumber)

	toUser, err := h.GQL.GetOrCreateUser(ctx, event.To.Hex())
	if err != nil {
		return fmt.Errorf("failed to ensure 'to' user exists: %w", err)
	}
	// h.AddOperation("User", toUser, event.Raw.BlockHash.Hex(), event.Raw.BlockNumber)

	tokenID := event.Id.String()

	// Create or update item
	item, err := h.GQL.GetItemByTokenId(ctx, tokenID)
	if err != nil {
		item, err = h.GQL.CreateItem(ctx, graphqldb.CreateItemParams{
			ID:       uuid.New().String(),
			TokenID:  tokenID,
			TokenUri: "",
			Standard: "ERC1155",
		})
		if err != nil {
			return fmt.Errorf("failed to create item: %w", err)
		}
	}

	h.Logger.Infow("TransferSingle event details",
		"from", event.From.Hex(),
		"to", event.To.Hex(),
		"operator", event.Operator.Hex(),
		"id", event.Id.String(),
		"value", event.Value.String(),
		"block_number", event.Raw.BlockNumber,
		"tx_hash", event.Raw.TxHash.Hex(),
		"contract", event.Raw.Address.Hex(),
	)

	if event.From.Hex() != "0x0000000000000000000000000000000000000000" {
		// Get current balance and subtract transferred amount
		fromBalance, err := h.GQL.GetUserBalance(ctx, graphqldb.GetUserBalanceParams{
			OwnerID: fromUser.ID,
			ItemID:  item.ID,
		})
		if err != nil {
			return fmt.Errorf("failed to get sender balance: %w", err)
		}

		currentValue, _ := new(big.Int).SetString(fromBalance.Value, 10)
		if currentValue == nil {
			return fmt.Errorf("invalid balance value for sender: %s", fromBalance.Value)
		}

		newValue := new(big.Int).Sub(currentValue, event.Value)
		if newValue.Sign() < 0 {
			return fmt.Errorf("insufficient balance for transfer")
		}

		// Update sender's balance
		_, err = h.GQL.UpsertBalance(ctx, graphqldb.UpsertBalanceParams{
			ID:        fromBalance.ID, // Use existing balance ID
			ItemID:    item.ID,
			OwnerID:   fromUser.ID,
			Value:     newValue.String(),
			UpdatedAt: fmt.Sprint(event.Raw.BlockNumber),
			Contract:  event.Raw.Address.Hex(),
		})
		if err != nil {
			return fmt.Errorf("failed to update sender balance: %w", err)
		}
	} else {
		balance, err := h.GQL.CreateBalance(ctx, graphqldb.CreateBalanceParams{
			ID:        uuid.New().String(),
			ItemID:    item.ID,
			OwnerID:   toUser.ID,
			Value:     event.Value.String(),
			UpdatedAt: fmt.Sprint(event.Raw.BlockNumber),
			Contract:  event.Raw.Address.Hex(),
		})
		if err != nil {
			return fmt.Errorf("failed to create balance: %w", err)
		}
		// This is a mint operation - no need to check sender balance
		h.Logger.Infow("Minting new tokens",
			"to", event.To.Hex(),
			"tokenId", tokenID,
			"amount", event.Value.String(),
			"balanceId", balance,
		)
	}

	// Get current balance and add transferred amount
	toBalance, err := h.GQL.GetUserBalance(ctx, graphqldb.GetUserBalanceParams{
		OwnerID: toUser.ID,
		ItemID:  item.ID,
	})

	balanceID := uuid.New().String()
	currentValue := big.NewInt(0)

	if err == nil {
		// Use existing balance ID and value if found
		balanceID = toBalance.ID
		currentValue, _ = new(big.Int).SetString(toBalance.Value, 10)
		if currentValue == nil {
			currentValue = big.NewInt(0)
		}
	} else {
		// If no balance exists, start from 0
		h.Logger.Infow("Creating new balance record",
			"user", toUser.ID,
			"tokenId", tokenID,
		)
	}

	newValue := new(big.Int).Add(currentValue, event.Value)

	// Update receiver's balance
	_, err = h.GQL.UpsertBalance(ctx, graphqldb.UpsertBalanceParams{
		ID:        balanceID,
		ItemID:    item.ID,
		OwnerID:   toUser.ID,
		Value:     newValue.String(),
		UpdatedAt: fmt.Sprint(event.Raw.BlockNumber),
		Contract:  event.Raw.Address.Hex(),
	})
	if err != nil {
		return fmt.Errorf("failed to update 1155 receiver balance: %w", err)
	}

	if err := h.SubmitToDA(); err != nil {
		h.Logger.Errorw("Failed to submit to DA", "error", err)
	}

	return nil
}

func (h *TransferBatchHandler) HandleTransferBatch(ctx context.Context, event *eventhandlers.TransferBatch) error {
	fromUser, err := h.GQL.GetOrCreateUser(ctx, event.From.Hex())
	if err != nil {
		return fmt.Errorf("failed to ensure 'from' user exists: %w", err)
	}
	// h.AddOperation("User", fromUser, event.Raw.BlockHash.Hex(), event.Raw.BlockNumber)

	toUser, err := h.GQL.GetOrCreateUser(ctx, event.To.Hex())
	if err != nil {
		return fmt.Errorf("failed to ensure 'to' user exists: %w", err)
	}
	// h.AddOperation("User", toUser, event.Raw.BlockHash.Hex(), event.Raw.BlockNumber)

	// Handle each token transfer in the batch
	for i := range event.Ids {
		tokenID := event.Ids[i].String()
		value := event.Values[i]

		// Create or update item
		item, err := h.GQL.GetItemByTokenId(ctx, tokenID)
		if err != nil {
			item, err = h.GQL.CreateItem(ctx, graphqldb.CreateItemParams{
				ID:       uuid.New().String(),
				TokenID:  tokenID,
				TokenUri: "",
				Standard: "ERC1155",
			})
			if err != nil {
				return fmt.Errorf("failed to create item: %w", err)
			}
		}

		// Update balances similar to TransferSingle
		if event.From.Hex() != "0x0000000000000000000000000000000000000000" {
			fromBalance, err := h.GQL.GetUserBalance(ctx, graphqldb.GetUserBalanceParams{
				OwnerID: fromUser.ID,
				ItemID:  item.ID,
			})
			if err != nil {
				return fmt.Errorf("failed to get sender balance: %w", err)
			}

			currentValue, _ := new(big.Int).SetString(fromBalance.Value, 10)
			if currentValue == nil {
				return fmt.Errorf("invalid balance value for sender: %s", fromBalance.Value)
			}

			newValue := new(big.Int).Sub(currentValue, value)
			if newValue.Sign() < 0 {
				return fmt.Errorf("insufficient balance for transfer")
			}

			_, err = h.GQL.UpsertBalance(ctx, graphqldb.UpsertBalanceParams{
				ID:        fromBalance.ID,
				ItemID:    item.ID,
				OwnerID:   fromUser.ID,
				Value:     newValue.String(),
				UpdatedAt: fmt.Sprint(event.Raw.BlockNumber),
				Contract:  event.Raw.Address.Hex(),
			})
			if err != nil {
				return fmt.Errorf("failed to update sender balance: %w", err)
			}
		}

		toBalance, err := h.GQL.GetUserBalance(ctx, graphqldb.GetUserBalanceParams{
			OwnerID: toUser.ID,
			ItemID:  item.ID,
		})
		if err != nil {
			toBalance = graphqldb.GetUserBalanceRow{
				ID:    uuid.New().String(),
				Value: "0",
			}
		}

		currentValue, _ := new(big.Int).SetString(toBalance.Value, 10)
		if currentValue == nil {
			return fmt.Errorf("invalid balance value for receiver: %s", toBalance.Value)
		}

		newValue := new(big.Int).Add(currentValue, value)

		_, err = h.GQL.UpsertBalance(ctx, graphqldb.UpsertBalanceParams{
			ID:        toBalance.ID,
			ItemID:    item.ID,
			OwnerID:   toUser.ID,
			Value:     newValue.String(),
			UpdatedAt: fmt.Sprint(event.Raw.BlockNumber),
			Contract:  event.Raw.Address.Hex(),
		})
		if err != nil {
			return fmt.Errorf("failed to update receiver balance: %w", err)
		}
	}

	if err := h.SubmitToDA(); err != nil {
		h.Logger.Errorw("Failed to submit to DA", "error", err)
	}

	return nil
}

// func (h *URIHandler) HandleURI(ctx context.Context, event *eventhandlers.URI) error {
// 	tokenID := event.Id.String()

// 	// Update item's URI
// 	item, err := h.GQL.GetItemByTokenId(ctx, tokenID)
// 	if err == nil {
// 		item, err = h.GQL.UpdateItem(ctx, graphqldb.UpdateItemParams{
// 			ID:       item.ID,
// 			TokenID:  item.TokenID,
// 			TokenUri: event.Value,
// 			Standard: item.Standard,
// 		})
// 		if err != nil {
// 			return fmt.Errorf("failed to update item URI: %w", err)
// 		}
// 		h.AddOperation("Item", item, event.Raw.BlockHash.Hex(), event.Raw.BlockNumber)
// 	}

// 	if err := h.SubmitToDA(); err != nil {
// 		h.Logger.Errorw("Failed to submit to DA", "error", err)
// 	}

// 	return nil
// }

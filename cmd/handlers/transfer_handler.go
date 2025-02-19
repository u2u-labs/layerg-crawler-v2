package handlers

import (
	"context"
	"fmt"

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

func (h *TransferHandler) HandleTransfer(ctx context.Context, event *eventhandlers.Transfer) error {
	fromUser, err := h.GQL.GetOrCreateUser(ctx, event.From.Hex())
	if err != nil {
		h.Logger.Errorw("Failed to ensure 'from' user exists",
			"err", err,
			"address", event.From.Hex(),
		)
		return fmt.Errorf("failed to ensure 'from' user exists: %w", err)
	}
	h.AddOperation("User", fromUser, event.Raw.BlockHash.Hex(), event.Raw.BlockNumber)

	toUser, err := h.GQL.GetOrCreateUser(ctx, event.To.Hex())
	if err != nil {
		h.Logger.Errorw("Failed to ensure 'to' user exists",
			"err", err,
			"address", event.To.Hex(),
		)
		return fmt.Errorf("failed to ensure 'to' user exists: %w", err)
	}
	h.AddOperation("User", toUser, event.Raw.BlockHash.Hex(), event.Raw.BlockNumber)

	// Try to get existing item
	tokenID := event.TokenId.String()
	items, err := h.GQL.ListItem(ctx)
	if err != nil {
		return fmt.Errorf("failed to list items: %w", err)
	}

	var existingItem graphqldb.Item
	var itemExists bool
	for _, item := range items {
		if item.TokenID == tokenID {
			existingItem = item
			itemExists = true
			break
		}
	}

	if itemExists {
		// Update existing item's owner
		updatedItem, err := h.GQL.UpdateItem(ctx, graphqldb.UpdateItemParams{
			ID:       existingItem.ID,
			TokenID:  existingItem.TokenID,
			TokenUri: existingItem.TokenUri,
			OwnerID:  toUser.ID,
		})
		if err != nil {
			h.Logger.Errorw("Failed to update item ownership",
				"err", err,
				"tokenId", tokenID,
				"newOwner", toUser.ID,
			)
			return fmt.Errorf("failed to update item ownership: %w", err)
		}
		h.AddOperation("Item", updatedItem, event.Raw.BlockHash.Hex(), event.Raw.BlockNumber)
	} else {
		// Create new item
		newItem, err := h.GQL.CreateItem(ctx, graphqldb.CreateItemParams{
			ID:       uuid.New().String(),
			TokenID:  tokenID,
			TokenUri: "",
			OwnerID:  toUser.ID,
		})
		if err != nil {
			h.Logger.Errorw("Failed to create new item",
				"err", err,
				"tokenId", tokenID,
				"owner", toUser.ID,
			)
			return fmt.Errorf("failed to create new item: %w", err)
		}
		h.AddOperation("Item", newItem, event.Raw.BlockHash.Hex(), event.Raw.BlockNumber)
	}

	h.Logger.Infow("Transfer event processed",
		"tokenId", tokenID,
		"from", fromUser.ID,
		"to", toUser.ID,
		"contract", event.Raw.Address.Hex(),
		"tx", event.Raw.TxHash.Hex(),
	)

	// Submit all changes to DA
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

	h.AddOperation("User", actorUser, event.Raw.BlockHash.Hex(), event.Raw.BlockNumber)

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

			itemResp, err := h.GQL.UpdateItem(ctx, graphqldb.UpdateItemParams{
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
			h.AddOperation("Item", itemResp, event.Raw.BlockHash.Hex(), event.Raw.BlockNumber)
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

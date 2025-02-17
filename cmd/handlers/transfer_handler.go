package handlers

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/google/uuid"
	graphqldb "github.com/u2u-labs/layerg-crawler/db/graphqldb"
	db "github.com/u2u-labs/layerg-crawler/db/sqlc"
	"github.com/u2u-labs/layerg-crawler/generated/eventhandlers"
	"github.com/u2u-labs/layerg-crawler/generated/mappings"
)

type TransferHandler struct {
	*mappings.BaseTransferMapping
}

type MetadataUpdateHandler struct {
	*mappings.BaseMetadataUpdateMapping
}

func NewTransferHandler(queries *db.Queries, gqlQueries *graphqldb.Queries, chainID int32, logger *zap.SugaredLogger) *TransferHandler {
	return &TransferHandler{
		BaseTransferMapping: mappings.NewTransferMapping(queries, gqlQueries, chainID, logger),
	}
}

func NewMetadataUpdateHandler(queries *db.Queries, gqlQueries *graphqldb.Queries, chainID int32, logger *zap.SugaredLogger) *MetadataUpdateHandler {
	return &MetadataUpdateHandler{
		BaseMetadataUpdateMapping: mappings.NewMetadataUpdateMapping(queries, gqlQueries, chainID, logger),
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

	toUser, err := h.GQL.GetOrCreateUser(ctx, event.To.Hex())
	if err != nil {
		h.Logger.Errorw("Failed to ensure 'to' user exists",
			"err", err,
			"address", event.To.Hex(),
		)
		return fmt.Errorf("failed to ensure 'to' user exists: %w", err)
	}

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
		_, err = h.GQL.UpdateItem(ctx, graphqldb.UpdateItemParams{
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
	} else {
		// Create new item
		_, err = h.GQL.CreateItem(ctx, graphqldb.CreateItemParams{
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
	}

	h.Logger.Infow("Transfer event processed",
		"tokenId", tokenID,
		"from", fromUser.ID,
		"to", toUser.ID,
		"contract", event.Raw.Address.Hex(),
		"tx", event.Raw.TxHash.Hex(),
	)

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

			_, err = h.GQL.UpdateItem(ctx, graphqldb.UpdateItemParams{
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
			break
		}
	}

	h.Logger.Infow("Metadata update event processed",
		"tokenId", tokenID,
		"actor", actor,
		"contract", event.Raw.Address.Hex(),
		"tx", event.Raw.TxHash.Hex(),
	)

	return nil
}

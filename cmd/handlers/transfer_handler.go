package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/google/uuid"
	graphqldb "github.com/u2u-labs/layerg-crawler/db/graphqldb"
	db "github.com/u2u-labs/layerg-crawler/db/sqlc"
	"github.com/u2u-labs/layerg-crawler/generated/eventhandlers"
	"github.com/u2u-labs/layerg-crawler/generated/mappings"
)

// TransferHandler implements mappings.TransferHandler
type TransferHandler struct {
	*mappings.BaseTransferMapping
}
type ApprovalHandler struct {
	*mappings.BaseApprovalMapping
}

func NewTransferHandler(queries *db.Queries, gqlQueries *graphqldb.Queries, chainID int32, logger *zap.SugaredLogger) *TransferHandler {
	return &TransferHandler{
		BaseTransferMapping: mappings.NewTransferMapping(queries, gqlQueries, chainID, logger),
	}
}
func NewApprovalHandler(queries *db.Queries, gqlQueries *graphqldb.Queries, chainID int32, logger *zap.SugaredLogger) *ApprovalHandler {
	return &ApprovalHandler{
		BaseApprovalMapping: mappings.NewApprovalMapping(queries, gqlQueries, chainID, logger),
	}
}

func (h *ApprovalHandler) HandleApproval(ctx context.Context, event *eventhandlers.Approval) error {
	return nil
}

// HandleTransfer implements mappings.TransferHandler
func (h *TransferHandler) HandleTransfer(ctx context.Context, event *eventhandlers.Transfer) error {
	// Store the transfer in system database
	// _, err := h.Queries.CreateOnchainHistory(ctx, db.CreateOnchainHistoryParams{
	// 	From:      event.From.Hex(),
	// 	To:        event.To.Hex(),
	// 	ChainID:   h.ChainID,
	// 	AssetID:   event.Raw.Address.Hex(),
	// 	TxHash:    event.Raw.TxHash.Hex(),
	// 	Receipt:   []byte("{}"),
	// 	EventType: sql.NullString{String: "Transfer", Valid: true},
	// 	Timestamp: time.Now(),
	// })

	// if err != nil {
	// 	h.Logger.Errorw("Failed to store transfer",
	// 		"chain_id", h.ChainID,
	// 		"err", err,
	// 		"tx", event.Raw.TxHash.Hex(),
	// 	)
	// 	return err
	// }

	// Create a Collection record in the GraphQL database
	_, err := h.GQL.CreateTransfer(ctx, graphqldb.CreateTransferParams{
		ID:        uuid.New().String(),
		From:      event.From.Hex(),
		To:        event.To.Hex(),
		Amount:    sql.NullString{String: event.Value.String(), Valid: true},
		Timestamp: sql.NullTime{Time: time.Now(), Valid: true},
	})

	if err != nil {
		h.Logger.Errorw("Failed to create transfer record",
			"err", err,
			"address", event.Raw.Address.Hex(),
			"from", event.From.Hex(),
			"to", event.To.Hex(),
			"amount", event.Value.String(),
		)
		return fmt.Errorf("failed to create transfer record: %w", err)
	}

	h.Logger.Infow("Transfer event processed",
		"from", event.From.Hex(),
		"to", event.To.Hex(),
		"amount", event.Value.String(),
		"contract", event.Raw.Address.Hex(),
		"tx", event.Raw.TxHash.Hex(),
	)

	return nil
}

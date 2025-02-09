package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/unicornultrafoundation/go-u2u/core/types"
	"go.uber.org/zap"

	graphqldb "github.com/u2u-labs/layerg-crawler/db/graphqldb"
	db "github.com/u2u-labs/layerg-crawler/db/sqlc"
	"github.com/u2u-labs/layerg-crawler/generated/eventhandlers"
)

type TransferHandler struct {
	queries    *db.Queries
	gqlQueries *graphqldb.Queries
	chainID    int32
}

func NewTransferHandler(queries *db.Queries, gqlQueries *graphqldb.Queries, chainID int32) *TransferHandler {
	return &TransferHandler{
		queries:    queries,
		gqlQueries: gqlQueries,
		chainID:    chainID,
	}
}

func (h *TransferHandler) HandleEvent(ctx context.Context, log *types.Log, logger *zap.SugaredLogger) error {
	event, err := eventhandlers.UnpackTransfer(log)
	if err != nil {
		return fmt.Errorf("failed to unpack event: %w", err)
	}

	// Store the transfer in system database
	_, err = h.queries.CreateOnchainHistory(ctx, db.CreateOnchainHistoryParams{
		From:      event.From.Hex(),
		To:        event.To.Hex(),
		ChainID:   h.chainID,
		AssetID:   log.Address.Hex(),
		TxHash:    log.TxHash.Hex(),
		Receipt:   []byte("{}"),
		EventType: sql.NullString{String: "Transfer", Valid: true},
		Timestamp: time.Now(),
	})

	if err != nil {
		logger.Errorw("Failed to store transfer",
			"err", err,
			"tx", log.TxHash.Hex(),
		)
		return err
	}

	// Create a Collection record in the GraphQL database
	_, err = h.gqlQueries.CreateCollection(ctx, graphqldb.CreateCollectionParams{
		ID:      uuid.New().String(),
		Address: log.Address.Hex(),
		Type: sql.NullString{
			String: "ERC20", // Or determine based on contract type
			Valid:  true,
		},
	})

	if err != nil {
		logger.Errorw("Failed to create collection",
			"err", err,
			"address", log.Address.Hex(),
		)
		return err
	}

	logger.Infow("Transfer event processed",
		"from", event.From.Hex(),
		"to", event.To.Hex(),
		"amount", event.Amount.String(),
		"contract", log.Address.Hex(),
		"tx", log.TxHash.Hex(),
	)

	return nil
}

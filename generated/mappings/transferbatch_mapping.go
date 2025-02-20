// Code generated - DO NOT EDIT.
package mappings

import (
	"context"

	"go.uber.org/zap"

	"github.com/u2u-labs/layerg-crawler/generated/eventhandlers"
	"github.com/u2u-labs/layerg-crawler/db/graphqldb"
	dbCon "github.com/u2u-labs/layerg-crawler/db/sqlc"
)

// TransferBatchHandler defines the interface for handling TransferBatch events
type TransferBatchHandler interface {
	HandleTransferBatch(ctx context.Context, event *eventhandlers.TransferBatch) error
}

// BaseTransferBatchMapping provides the base structure for TransferBatch event mappings
type BaseTransferBatchMapping struct {
	Queries  *dbCon.Queries
	GQL      *graphqldb.Queries
	ChainID  int32
	Logger   *zap.SugaredLogger
}

// NewTransferBatchMapping creates a new base mapping for TransferBatch events
func NewTransferBatchMapping(queries *dbCon.Queries, gql *graphqldb.Queries, chainID int32, logger *zap.SugaredLogger) *BaseTransferBatchMapping {
	return &BaseTransferBatchMapping{
		Queries:  queries,
		GQL:      gql,
		ChainID:  chainID,
		Logger:   logger,
	}
}

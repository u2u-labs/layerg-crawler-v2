// Code generated - DO NOT EDIT.
package mappings

import (
	"context"

	"go.uber.org/zap"

	"github.com/u2u-labs/layerg-crawler/generated/eventhandlers"
	"github.com/u2u-labs/layerg-crawler/db/graphqldb"
	dbCon "github.com/u2u-labs/layerg-crawler/db/sqlc"
)

// TransferHandler defines the interface for handling Transfer events
type TransferHandler interface {
	HandleTransfer(ctx context.Context, event *eventhandlers.Transfer) error
}

// BaseTransferMapping provides the base structure for Transfer event mappings
type BaseTransferMapping struct {
	Queries  *dbCon.Queries
	GQL      *graphqldb.Queries
	ChainID  int32
	Logger   *zap.SugaredLogger
}

// NewTransferMapping creates a new base mapping for Transfer events
func NewTransferMapping(queries *dbCon.Queries, gql *graphqldb.Queries, chainID int32, logger *zap.SugaredLogger) *BaseTransferMapping {
	return &BaseTransferMapping{
		Queries:  queries,
		GQL:      gql,
		ChainID:  chainID,
		Logger:   logger,
	}
}

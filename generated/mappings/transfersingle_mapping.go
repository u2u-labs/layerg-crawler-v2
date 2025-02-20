// Code generated - DO NOT EDIT.
package mappings

import (
	"context"

	"go.uber.org/zap"

	"github.com/u2u-labs/layerg-crawler/generated/eventhandlers"
	"github.com/u2u-labs/layerg-crawler/db/graphqldb"
	dbCon "github.com/u2u-labs/layerg-crawler/db/sqlc"
)

// TransferSingleHandler defines the interface for handling TransferSingle events
type TransferSingleHandler interface {
	HandleTransferSingle(ctx context.Context, event *eventhandlers.TransferSingle) error
}

// BaseTransferSingleMapping provides the base structure for TransferSingle event mappings
type BaseTransferSingleMapping struct {
	Queries  *dbCon.Queries
	GQL      *graphqldb.Queries
	ChainID  int32
	Logger   *zap.SugaredLogger
}

// NewTransferSingleMapping creates a new base mapping for TransferSingle events
func NewTransferSingleMapping(queries *dbCon.Queries, gql *graphqldb.Queries, chainID int32, logger *zap.SugaredLogger) *BaseTransferSingleMapping {
	return &BaseTransferSingleMapping{
		Queries:  queries,
		GQL:      gql,
		ChainID:  chainID,
		Logger:   logger,
	}
}

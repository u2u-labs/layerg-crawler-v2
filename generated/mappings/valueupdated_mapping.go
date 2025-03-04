// Code generated - DO NOT EDIT.
package mappings

import (
	"context"

	"go.uber.org/zap"

	"github.com/u2u-labs/layerg-crawler/generated/eventhandlers"
	"github.com/u2u-labs/layerg-crawler/db/graphqldb"
	dbCon "github.com/u2u-labs/layerg-crawler/db/sqlc"
)

// ValueUpdatedHandler defines the interface for handling ValueUpdated events
type ValueUpdatedHandler interface {
	HandleValueUpdated(ctx context.Context, event *eventhandlers.ValueUpdated) error
}

// BaseValueUpdatedMapping provides the base structure for ValueUpdated event mappings
type BaseValueUpdatedMapping struct {
	Queries  *dbCon.Queries
	GQL      *graphqldb.Queries
	ChainID  int32
	Logger   *zap.SugaredLogger
}

// NewValueUpdatedMapping creates a new base mapping for ValueUpdated events
func NewValueUpdatedMapping(queries *dbCon.Queries, gql *graphqldb.Queries, chainID int32, logger *zap.SugaredLogger) *BaseValueUpdatedMapping {
	return &BaseValueUpdatedMapping{
		Queries:  queries,
		GQL:      gql,
		ChainID:  chainID,
		Logger:   logger,
	}
}

// Code generated - DO NOT EDIT.
package mappings

import (
	"context"

	"go.uber.org/zap"

	"github.com/u2u-labs/layerg-crawler/generated/eventhandlers"
	"github.com/u2u-labs/layerg-crawler/db/graphqldb"
	dbCon "github.com/u2u-labs/layerg-crawler/db/sqlc"
)

// MetadataUpdateHandler defines the interface for handling MetadataUpdate events
type MetadataUpdateHandler interface {
	HandleMetadataUpdate(ctx context.Context, event *eventhandlers.MetadataUpdate) error
}

// BaseMetadataUpdateMapping provides the base structure for MetadataUpdate event mappings
type BaseMetadataUpdateMapping struct {
	Queries  *dbCon.Queries
	GQL      *graphqldb.Queries
	ChainID  int32
	Logger   *zap.SugaredLogger
}

// NewMetadataUpdateMapping creates a new base mapping for MetadataUpdate events
func NewMetadataUpdateMapping(queries *dbCon.Queries, gql *graphqldb.Queries, chainID int32, logger *zap.SugaredLogger) *BaseMetadataUpdateMapping {
	return &BaseMetadataUpdateMapping{
		Queries:  queries,
		GQL:      gql,
		ChainID:  chainID,
		Logger:   logger,
	}
}

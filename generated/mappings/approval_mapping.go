// Code generated - DO NOT EDIT.
package mappings

import (
	"context"

	"go.uber.org/zap"

	"github.com/u2u-labs/layerg-crawler/generated/eventhandlers"
	"github.com/u2u-labs/layerg-crawler/db/graphqldb"
	dbCon "github.com/u2u-labs/layerg-crawler/db/sqlc"
)

// ApprovalHandler defines the interface for handling Approval events
type ApprovalHandler interface {
	HandleApproval(ctx context.Context, event *eventhandlers.Approval) error
}

// BaseApprovalMapping provides the base structure for Approval event mappings
type BaseApprovalMapping struct {
	Queries  *dbCon.Queries
	GQL      *graphqldb.Queries
	ChainID  int32
	Logger   *zap.SugaredLogger
}

// NewApprovalMapping creates a new base mapping for Approval events
func NewApprovalMapping(queries *dbCon.Queries, gql *graphqldb.Queries, chainID int32, logger *zap.SugaredLogger) *BaseApprovalMapping {
	return &BaseApprovalMapping{
		Queries:  queries,
		GQL:      gql,
		ChainID:  chainID,
		Logger:   logger,
	}
}

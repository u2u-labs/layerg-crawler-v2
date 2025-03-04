// Code generated - DO NOT EDIT.
package router

import (
    "context"
    "fmt"

    "github.com/unicornultrafoundation/go-u2u/core/types"
    "go.uber.org/zap"

    "github.com/u2u-labs/layerg-crawler/generated/eventhandlers"
    "github.com/u2u-labs/layerg-crawler/db/graphqldb"
    dbCon "github.com/u2u-labs/layerg-crawler/db/sqlc"
    "github.com/u2u-labs/layerg-crawler/cmd/handlers"
)

type EventRouter struct {
    handlers map[string]interface{}
    queries  *dbCon.Queries
    gql      *graphqldb.Queries
    logger   *zap.SugaredLogger
    chainID  int32
}

func NewEventRouter(queries *dbCon.Queries, gql *graphqldb.Queries, logger *zap.SugaredLogger, chainID int32) *EventRouter {
    return &EventRouter{
        handlers: make(map[string]interface{}),
        queries:  queries,
        gql:      gql,
        logger:   logger,
        chainID:  chainID,
    }
}

func (r *EventRouter) Route(ctx context.Context, log *types.Log) error {
    if len(log.Topics) == 0 {
        return nil
    }

    signature := log.Topics[0].Hex()
    switch signature {
    case eventhandlers.ValueUpdatedEventSignature:
        handler := handlers.NewValueUpdatedHandler(r.queries, r.gql, r.chainID, r.logger)
        event, err := eventhandlers.UnpackValueUpdated(log)
        if err != nil {
            return fmt.Errorf("failed to unpack ValueUpdated event: %w", err)
        }
        return handler.HandleValueUpdated(ctx, event)
    default:
        r.logger.Debugw("No handler for event signature", "signature", signature)
        return nil
    }
}

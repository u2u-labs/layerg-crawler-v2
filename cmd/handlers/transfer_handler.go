package handlers

import (
	"context"

	"go.uber.org/zap"

	"github.com/u2u-labs/layerg-crawler/cmd/helpers"
	graphqldb "github.com/u2u-labs/layerg-crawler/db/graphqldb"
	db "github.com/u2u-labs/layerg-crawler/db/sqlc"
	"github.com/u2u-labs/layerg-crawler/generated/abi_helpers"
	"github.com/u2u-labs/layerg-crawler/generated/eventhandlers"
)

type HandleValueUpdated struct {
	*BaseHandler
}

func NewValueUpdatedHandler(queries *db.Queries, gqlQueries *graphqldb.Queries, chainID int32, logger *zap.SugaredLogger) *HandleValueUpdated {
	return &HandleValueUpdated{
		BaseHandler: &BaseHandler{
			Queries: queries,
			GQL:     gqlQueries,
			ChainID: chainID,
			Logger:  logger,
		},
	}
}

func (h *HandleValueUpdated) HandleValueUpdated(ctx context.Context, event *eventhandlers.ValueUpdated) error {

	client, err := helpers.InitChainClient("https://rpc-nebulas-testnet.uniultra.xyz")
	if err != nil {
		return err
	}

	value, err := abi_helpers.SIMPLESTORAGEGetValue(ctx, h.Logger, client, &event.Raw.Address)
	h.Logger.Info("GetValue", "hash", value)

	if err != nil {
		return err
	}

	sender, err := abi_helpers.SIMPLESTORAGESender(ctx, h.Logger, client, &event.Raw.Address)
	h.Logger.Info("Sender", "hash", sender)

	if err != nil {
		return err
	}

	h.GQL.CreateValue(ctx, graphqldb.CreateValueParams{
		ID:     event.Raw.TxHash.Hex(),
		Value:  value.String(),
		Sender: sender.String(),
	})

	return nil
}

package controllers

import (
	"context"
	"database/sql"
	"math/big"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	db "github.com/u2u-labs/layerg-crawler/db/sqlc"
	u2u "github.com/unicornultrafoundation/go-u2u"
	"github.com/unicornultrafoundation/go-u2u/common"
	"github.com/unicornultrafoundation/go-u2u/ethclient"
)

type BackFillController struct {
	db    *db.Queries
	rawDb *sql.DB
	ctx   context.Context
	rdb   *redis.Client
}

func NewBackFillController(db *db.Queries, rawDb *sql.DB, ctx context.Context, rdb *redis.Client) *BackFillController {
	return &BackFillController{db, rawDb, ctx, rdb}
}

func (bfc *BackFillController) BackFill(ctx *gin.Context) {
	chainId := ctx.Param("chain_id")
	contractAddress := ctx.Param("contract_address")
	startingBlockStr := ctx.Query("starting_block")
	startingBlock, err := strconv.Atoi(startingBlockStr)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid starting_block"})
		return
	}

	chainIdInt, err := strconv.Atoi(chainId)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid chain_id"})
		return
	}
	chain, err := bfc.db.GetChainById(bfc.ctx, int32(chainIdInt))

	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid chain_id"})
	}

	client, err := initChainClient(&chain)
	if err != nil {
		return
	}

	logs, _ := client.FilterLogs(ctx, u2u.FilterQuery{
		BlockHash: nil,
		FromBlock: big.NewInt(int64(startingBlock)),
		ToBlock:   big.NewInt(25574300),
		Addresses: []common.Address{common.HexToAddress(contractAddress)},
	})

	// sample := logs[1]

	// fmt.Println(sample.Topics[0].Hex())

	ctx.JSON(200, gin.H{"message": "Backfilling started", "logs": logs})

}

func initChainClient(chain *db.Chain) (*ethclient.Client, error) {
	return ethclient.Dial(chain.RpcUrl)
}

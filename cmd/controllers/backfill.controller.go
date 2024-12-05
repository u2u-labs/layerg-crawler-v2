package controllers

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/u2u-labs/layerg-crawler/cmd/response"
	db "github.com/u2u-labs/layerg-crawler/db/sqlc"
	"github.com/unicornultrafoundation/go-u2u/common"
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

func (bfc *BackFillController) AddBackFillTracker(ctx *gin.Context) {
	var params *db.AddBackfillCrawlerParams

	// Read the request body
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ErrorResponseData(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// convert params.contractAddress to checksum
	params.CollectionAddress = common.HexToAddress(params.CollectionAddress).Hex()

	// add to db
	if err := bfc.db.AddBackfillCrawler(ctx, *params); err != nil {
		response.ErrorResponseData(ctx, http.StatusInternalServerError, err.Error())
		return
	}
}

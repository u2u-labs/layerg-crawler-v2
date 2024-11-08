package controllers

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/u2u-labs/layerg-crawler/cmd/response"
	db "github.com/u2u-labs/layerg-crawler/db/sqlc"
)

type HistoryController struct {
	db    *db.Queries
	rawDb *sql.DB
	ctx   context.Context
	rdb   *redis.Client
}

func NewHistoryController(db *db.Queries, rawDb *sql.DB, ctx context.Context, rdb *redis.Client) *HistoryController {
	return &HistoryController{db, rawDb, ctx, rdb}
}

// Get onchain history godoc
// @Summary      Get History of a transaction
// @Description  Get History of a transaction
// @Tags         history
// @Accept       json
// @Produce      json
// @Param tx_hash query string true "Tx Hash"
// @Success      200 {object} response.ResponseData
// @Security     ApiKeyAuth
// @Router       /history [get]
func (hs *HistoryController) GetHistory(ctx *gin.Context) {
	txHash := ctx.Query("tx_hash")

	history, err := hs.db.GetOnchainHistoryByTxHash(ctx, txHash)
	if err != nil {
		response.ErrorResponseData(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessReponseData(ctx, http.StatusOK, history)
}

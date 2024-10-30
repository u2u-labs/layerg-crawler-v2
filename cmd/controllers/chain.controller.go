package controllers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/u2u-labs/layerg-crawler/cmd/response"
	db "github.com/u2u-labs/layerg-crawler/db/sqlc"
)

type ChainController struct {
	db  *db.Queries
	ctx context.Context
}

func NewChainController(db *db.Queries, ctx context.Context) *ChainController {
	return &ChainController{db, ctx}
}

// AddNewChain godoc
// @Summary      Add a new chain
// @Description  Add a new chain
// @Tags         chains
// @Accept       json
// @Produce      json
// @Param body body db.AddChainParams true "Add a new chain"
// @Security     BasicAuth
// @Success      200 {object} response.ResponseData
// @Failure      500 {object} response.ErrorResponse
// @Example      { "id": 1, "chain": "U2U", "name": "Nebulas Testnet", "RpcUrl": "sre", "ChainId": 2484, "Explorer": "str", "BlockTime": 500 }
// @Router       /chain [post]
func (cc *ChainController) AddNewChain(ctx *gin.Context) {
	var params *db.AddChainParams

	// Read the request body
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ErrorResponseData(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// add to db
	if err := cc.db.AddChain(ctx, *params); err != nil {
		response.ErrorResponseData(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// response
	resData, err := cc.db.GetChainById(ctx, params.ID)
	if err != nil {
		response.ErrorResponseData(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessReponseData(ctx, http.StatusOK, resData)
}

// GetAllSupportedChain godoc
// @Summary      Get all supported chains
// @Description  Get all supported chains
// @Tags         chains
// @Accept       json
// @Produce      json
// @Security     BasicAuth
// @Param chain_id query string false "Chain ID"
// @Router       /chain [get]
func (cc *ChainController) GetAllChains(ctx *gin.Context) {

	chainIdStr := ctx.Query("chain_id")

	if chainIdStr != "" {
		// get chain by id
		chainId, err := strconv.Atoi(chainIdStr)
		if err != nil {
			response.ErrorResponseData(ctx, http.StatusInternalServerError, err.Error())
			return
		}
		chain, err := cc.db.GetChainById(ctx, int32(chainId))
		if err != nil {
			response.ErrorResponseData(ctx, http.StatusInternalServerError, err.Error())
			return
		}
		response.SuccessReponseData(ctx, http.StatusOK, chain)
		return
	}

	// get all chains
	chains, err := cc.db.GetAllChain(ctx)
	if err != nil {
		response.ErrorResponseData(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.SuccessReponseData(ctx, http.StatusOK, chains)

}

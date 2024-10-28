package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
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
// @Param body body db.AddChainParams true "Chain network information"
// @Example      { "id": 1, "chain": "U2U", "name": "Nebulas Testnet", "RpcUrl": "sre", "ChainId": 2484, "Explorer": "str", "BlockTime": 500 }
// @Router       /chain [post]
func (cc *ChainController) AddNewChain(ctx *gin.Context) {
	var params *db.AddChainParams

	// Read the request body
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// add to db
	if err := cc.db.AddChain(ctx, *params); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Chain added", "data": params})
}

// GetAllSupportedChain godoc
// @Summary      Get all supported chains
// @Description  Get all supported chains
// @Tags         chains
// @Accept       json
// @Produce      json
// @Router       /chain [get]
func (cc *ChainController) GetAllChains(ctx *gin.Context) {
	chains, err := cc.db.GetAllChain(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": chains})
}

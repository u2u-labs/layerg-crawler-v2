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

// Get a single handler
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

// Get all supported chains
func (cc *ChainController) GetAllChains(ctx *gin.Context) {
	chains, err := cc.db.GetAllChain(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": chains})
}

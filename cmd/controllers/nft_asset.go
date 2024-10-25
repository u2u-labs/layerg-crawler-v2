package controllers

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "github.com/u2u-labs/layerg-crawler/db/sqlc"
)

type NftAssetController struct {
	db  *db.Queries
	ctx context.Context
}

func NewNftAssetController(db *db.Queries, ctx context.Context) *NftAssetController {
	return &NftAssetController{db, ctx}
}

// Get a single handler
func (cc *NftAssetController) GetErc721AssetByChainIdAddress(ctx *gin.Context) {
	chainIdStr := ctx.Query("chainId")
	chainId, err := strconv.Atoi(chainIdStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Invalid chainId"})
		return
	}
	contractAddress := ctx.Param("contractAddress")

	args := &db.GetAssetByChainIdAndContractAddressParams{
		ChainID:           int32(chainId),
		CollectionAddress: contractAddress,
	}

	assetCollection, err := cc.db.GetAssetByChainIdAndContractAddress(ctx, *args)

	switch assetType := assetCollection.Type; assetType {
	case db.AssetTypeERC721:
		erc721Assets, _ := cc.db.Get721AssetByAssetId(ctx, assetCollection.ID)
		ctx.JSON(http.StatusOK, gin.H{"status": "Successfully retrived id", "type": "ERC721", "asset": erc721Assets})

	}
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Failed to retrieve NFT with this contract address in the chain"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed retrieving NFT Asset", "error": err.Error()})
		return
	}

}

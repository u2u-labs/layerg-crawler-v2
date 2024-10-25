package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "github.com/u2u-labs/layerg-crawler/db/sqlc"
)

type AssetController struct {
	db  *db.Queries
	ctx context.Context
}

func NewAssetController(db *db.Queries, ctx context.Context) *AssetController {
	return &AssetController{db, ctx}
}

// Create new Asset
func (cc *AssetController) AddNewAsset(ctx *gin.Context) {
	var params *db.AddNewAssetParams
	chainIdStr := ctx.Param("chainId")
	chainId, err := strconv.Atoi(chainIdStr)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Invalid chainId"})
		return
	}

	// Read the raw body
	rawBodyData, err := ctx.GetRawData()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read body"})
		return
	}

	// Unmarshal JSON into the struct
	if err := json.Unmarshal(rawBodyData, &params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params.ChainID = int32(chainId)
	params.ID = strconv.Itoa(int(chainId)) + ":" + params.CollectionAddress

	ctx.JSON(http.StatusOK, gin.H{"message": "Asset added", "data": params})
}

// Get a single handler
func (cc *AssetController) GetAssetByChainIdAndContractAddress(ctx *gin.Context) {
	chainIdStr := ctx.Query("chainId")
	chainId, err := strconv.Atoi(chainIdStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Invalid chainId"})
		return
	}
	contractAddress := ctx.Query("contractAddress")

	tokenId := ctx.Query("tokenId")

	args := &db.GetAssetByChainIdAndContractAddressParams{
		ChainID:           int32(chainId),
		CollectionAddress: contractAddress,
	}

	assetCollection, err := cc.db.GetAssetByChainIdAndContractAddress(ctx, *args)

	switch assetType := assetCollection.Type; assetType {
	case db.AssetTypeERC721:
		erc721Assets, _ := func() (interface{}, error) {
			if tokenId != "" {
				args := &db.Get721AssetByAssetIdAndTokenIdParams{
					AssetID: assetCollection.ID,
					TokenID: tokenId,
				}
				return cc.db.Get721AssetByAssetIdAndTokenId(ctx, *args)
			}
			return cc.db.Get721AssetByAssetId(ctx, assetCollection.ID)
		}()
		ctx.JSON(http.StatusOK, gin.H{"status": "Successfully retrived id", "type": "ERC721", "asset": erc721Assets})
	case db.AssetTypeERC1155:
		erc1155Assets, _ := func() (interface{}, error) {
			if tokenId != "" {
				args := &db.Get1155AssetByAssetIdAndTokenIdParams{
					AssetID: assetCollection.ID,
					TokenID: tokenId,
				}
				return cc.db.Get1155AssetByAssetIdAndTokenId(ctx, *args)
			}
			return cc.db.Get1155AssetByAssetId(ctx, assetCollection.ID)
		}()
		ctx.JSON(http.StatusOK, gin.H{"status": "Successfully retrived id", "type": "ERC1155", "asset": erc1155Assets})
	case db.AssetTypeERC20:
		erc20Assets, _ := func() (interface{}, error) {
			if tokenId != "" {
				// error because there is no token id in erc20
				return nil, nil
			}
			asset, err := cc.db.Get20AssetByAssetId(ctx, assetCollection.ID)
			return asset, err
		}()
		ctx.JSON(http.StatusOK, gin.H{"status": "Successfully retrived id", "type": "ERC20", "asset": erc20Assets})
	}
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Failed to retrieve NFT with this contract address in the chain"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed retrieving Asset", "error": err.Error()})
		return
	}
}

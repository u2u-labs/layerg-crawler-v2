package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/u2u-labs/layerg-crawler/cmd/utils"
	db "github.com/u2u-labs/layerg-crawler/db/sqlc"
)

type AssetController struct {
	db    *db.Queries
	rawDb *sql.DB
	ctx   context.Context
}

func NewAssetController(db *db.Queries, rawDb *sql.DB, ctx context.Context) *AssetController {
	return &AssetController{db, rawDb, ctx}
}

// AddNewAsset godoc
// @Summary      Add a new asset collection to the chain
// @Description  Add a new asset collection to the chain
// @Tags         asset
// @Accept       json
// @Produce      json
// @Param chainId path int true "Chain ID"
// @Param body body utils.AddNewAssetParamsSwagger true "Asset collection information"
// @Example      { "id": 1, "chain": "U2U", "name": "Nebulas Testnet", "RpcUrl": "sre", "ChainId": 2484, "Explorer": "str", "BlockTime": 500 }
// @Router       /chain/:chainId/asset [post]
func (cc *AssetController) AddNewAsset(ctx *gin.Context) {
	// var params *db.AddNewAssetParams
	var params *utils.AddNewAssetParamsUtil
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

	assetParam := utils.ConvertUtilToParams(params)

	// add to db
	if err := cc.db.AddNewAsset(ctx, assetParam); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Asset added", "data": params})
}

// AddNewAsset godoc
// @Summary      Get all asset collection of the chain
// @Description  Get all asset collection of the chain
// @Tags         asset
// @Accept       json
// @Produce      json
// @Param chainId path int true "Chain ID"
// @Router       /chain/:chainId/asset [get]
func (cc *AssetController) GetAssetByChainId(ctx *gin.Context) {
	chainIdStr := ctx.Param("chainId")
	chainId, err := strconv.Atoi(chainIdStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Invalid chainId"})
		return
	}

	page, limit, offset := db.GetLimitAndOffset(ctx)

	assets, err := cc.db.GetPaginatedAssetsByChainId(ctx, db.GetPaginatedAssetsByChainIdParams{
		ChainID: int32(chainId),
		Limit:   int32(limit),
		Offset:  int32(offset),
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Query total items count
	totalAssets, err := cc.db.CountAssetByChainId(ctx, int32(chainId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create pagination response
	paginationResponse := db.Pagination[db.Asset]{
		Page:       page,
		Limit:      limit,
		TotalItems: totalAssets,
		TotalPages: (totalAssets + int64(limit) - 1) / int64(limit), // Calculate total pages
		Data:       assets,
	}

	ctx.JSON(http.StatusOK, gin.H{"data": paginationResponse})
}

// func (cc *AssetController) Test(ctx *gin.Context) {
// 	fmt.Print("Test")
// 	chainIdStr := ctx.Param("chainId")
// 	chainId, err := strconv.Atoi(chainIdStr)
// 	// if err != nil {
// 	// 	ctx.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Invalid chainId"})
// 	// 	return
// 	// }

// 	fmt.Print(chainId)

// 	filterConditions := make(map[string]string)
// 	if filterField := ctx.Query("collection_address"); filterField != "" {
// 		filterConditions["collection_address"] = filterField
// 	}

// 	fmt.Print(filterConditions)

// 	countAssetByCollectionAddress, err := db.CountItemsWithFilter(cc.rawDb, "assets", filterConditions)

// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{"data": countAssetByCollectionAddress})
// }

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

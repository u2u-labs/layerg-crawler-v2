package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/u2u-labs/layerg-crawler/cmd/response"
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
// @Param chain_id path string true "Chain ID"
// @Security     BasicAuth
// @Param body body utils.AddNewAssetParamsSwagger true "Asset collection information"
// @Example      { "id": 1, "chain": "U2U", "name": "Nebulas Testnet", "RpcUrl": "sre", "ChainId": 2484, "Explorer": "str", "BlockTime": 500 }
// @Router       /chain/{chain_id}/collection [post]
func (cc *AssetController) AddNewAsset(ctx *gin.Context) {
	var params *utils.AddNewAssetParamsUtil
	chainIdStr := ctx.Param("chain_id")
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

	assetParam := utils.ConvertCustomTypeToSqlParams(params)

	// add to db
	if err := cc.db.AddNewAsset(ctx, assetParam); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Output the result
	jsonResponse, err := utils.MarshalAssetParams(assetParam)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response.SuccessReponseData(ctx, http.StatusCreated, jsonResponse)
}

func (cc *AssetController) GetAssetByChainId(ctx *gin.Context) {
	chainIdStr := ctx.Param("chain_id")

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
	paginationResponse := db.Pagination[utils.AssetResponse]{
		Page:       page,
		Limit:      limit,
		TotalItems: totalAssets,
		TotalPages: (totalAssets + int64(limit) - 1) / int64(limit), // Calculate total pages
		Data:       utils.ConvertToAssetResponses(assets),
	}

	response.SuccessReponseData(ctx, http.StatusOK, paginationResponse)
}

// GetAssetCollection godoc
// @Summary      Get all asset in a collection of the chain
// @Description  Retrieve all asset collections associated with the specified chain ID.
// @Tags         asset
// @Accept       json
// @Produce      json
// @Security     BasicAuth
// @Param chain_id path string true "Chain ID"
// @Param collection_address query string false "Collection Address"
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Router       /chain/{chain_id}/collection [get]
func (cc *AssetController) GetAssetCollection(ctx *gin.Context) {

	hasQuery, collectionAddress := utils.GetAssetCollectionFilterParam(ctx)
	if !hasQuery {
		cc.GetAssetByChainId(ctx)
	} else {
		cc.GetAssetCollectionByChainIdAndContractAddress(ctx, collectionAddress)
	}

}

func (cc *AssetController) GetAssetCollectionByChainIdAndContractAddress(ctx *gin.Context, collectionAddress string) {
	chainIdStr := ctx.Param("chain_id")

	assetId := chainIdStr + ":" + collectionAddress

	assetCollection, err := cc.db.GetAssetById(ctx, assetId)

	if err != nil {
		if err == sql.ErrNoRows {
			response.ErrorResponseData(ctx, http.StatusNotFound, "Failed to retrieve asset collection with this contract address in the chain")
			return
		}
		response.ErrorResponseData(ctx, http.StatusBadGateway, err.Error())
		return
	}

	response.SuccessReponseData(ctx, http.StatusOK, utils.ConvertAssetToAssetResponse(assetCollection))
}

// GetAssetCollectionAByChainIdAndContractAddress godoc
// @Summary      Get all asset collection of the chain
// @Description  Get all asset collection of the chain
// @Tags         asset
// @Accept       json
// @Produce      json
// @Security     BasicAuth
// @Param chain_id path int true "Chain ID"
// @Param collection_address path string true "Collection Address"
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Param token_id query []string false "Token IDs" collectionFormat(multi)
// @Param owner query string false "Owner Address"
// @Router       /chain/{chain_id}/collection/{collection_address}/assets [get]
func (cc *AssetController) GetAssetByChainIdAndContractAddress(ctx *gin.Context) {
	chainIdStr := ctx.Param("chain_id")
	collectionAddress := ctx.Param("collection_address")

	assetId := chainIdStr + ":" + collectionAddress

	hasFilterParam, tokenIds, owner := utils.GetAssetFilterParam(ctx)

	if !hasFilterParam {
		cc.GetAssetsFromCollection(ctx, assetId)
	} else {
		cc.GetAssetsFromCollectionWithFilter(ctx, assetId, tokenIds, owner)
	}

}

func (cc *AssetController) GetAssetsFromCollection(ctx *gin.Context, assetId string) {
	assetCollection, err := cc.db.GetAssetById(ctx, assetId)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Failed to retrieve asset collection with this contract address in the chain"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed retrieving Asset", "error": err.Error()})
	}
	page, limit, offset := db.GetLimitAndOffset(ctx)

	switch assetType := assetCollection.Type; assetType {
	case db.AssetTypeERC721:
		totalAssets, err := cc.db.Count721AssetByAssetId(ctx, assetCollection.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		assets, _ := cc.db.GetPaginated721AssetByAssetId(ctx, db.GetPaginated721AssetByAssetIdParams{
			AssetID: assetCollection.ID,
			Limit:   int32(limit),
			Offset:  int32(offset),
		})

		paginationResponse := db.Pagination[utils.Erc721CollectionAssetResponse]{
			Page:       page,
			Limit:      limit,
			TotalItems: totalAssets,
			TotalPages: (totalAssets + int64(limit) - 1) / int64(limit),
			Data:       utils.ConvertToErc721CollectionAssetResponses(assets),
		}

		ctx.JSON(http.StatusOK, gin.H{"status": "Successfully retrived id", "type": "ERC721", "asset": paginationResponse})
	case db.AssetTypeERC1155:
		totalAssets, err := cc.db.Count1155AssetByAssetId(ctx, assetCollection.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		assets, _ := cc.db.GetPaginated1155AssetByAssetId(ctx, db.GetPaginated1155AssetByAssetIdParams{
			AssetID: assetCollection.ID,
			Limit:   int32(limit),
			Offset:  int32(offset),
		})

		paginationResponse := db.Pagination[utils.Erc1155CollectionAssetResponse]{
			Page:       page,
			Limit:      limit,
			TotalItems: totalAssets,
			TotalPages: (totalAssets + int64(limit) - 1) / int64(limit),
			Data:       utils.ConvertToErc1155CollectionAssetResponses(assets),
		}

		ctx.JSON(http.StatusOK, gin.H{"status": "Successfully retrived id", "type": "ERC1155", "asset": paginationResponse})

	case db.AssetTypeERC20:
		totalAssets, err := cc.db.Count20AssetByAssetId(ctx, assetCollection.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		assets, _ := cc.db.GetPaginated20AssetByAssetId(ctx, db.GetPaginated20AssetByAssetIdParams{
			AssetID: assetCollection.ID,
			Limit:   int32(limit),
			Offset:  int32(offset),
		})

		paginationResponse := db.Pagination[db.Erc20CollectionAsset]{
			Page:       page,
			Limit:      limit,
			TotalItems: totalAssets,
			TotalPages: (totalAssets + int64(limit) - 1) / int64(limit),
			Data:       assets,
		}

		ctx.JSON(http.StatusOK, gin.H{"status": "Successfully retrived id", "type": "ERC20", "asset": paginationResponse})
	}

}

func (cc *AssetController) GetAssetsFromCollectionWithFilter(ctx *gin.Context, assetId string, tokenIds []string, owner string) {
	assetCollection, err := cc.db.GetAssetById(ctx, assetId)
	filterConditions := make(map[string][]string)

	if len(tokenIds) > 0 {
		filterConditions["token_id"] = tokenIds
	}

	if owner != "" {
		filterConditions["owner"] = []string{owner}
	}

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Failed to retrieve asset collection with this contract address in the chain"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed retrieving Asset", "error": err.Error()})
	}
	page, limit, offset := db.GetLimitAndOffset(ctx)

	switch assetType := assetCollection.Type; assetType {
	case db.AssetTypeERC721:
		totalAssets, err := db.CountItemsWithFilter(cc.rawDb, "erc_721_collection_assets", filterConditions)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		assets, _ := db.QueryWithDynamicFilter[db.Erc721CollectionAsset](cc.rawDb, "erc_721_collection_assets", limit, offset, filterConditions)

		paginationResponse := db.Pagination[utils.Erc721CollectionAssetResponse]{
			Page:       page,
			Limit:      limit,
			TotalItems: int64(totalAssets),
			TotalPages: int64(totalAssets+(limit)-1) / int64(limit),
			Data:       utils.ConvertToErc721CollectionAssetResponses(assets),
		}

		ctx.JSON(http.StatusOK, gin.H{"status": "Successfully retrived id", "type": "ERC721", "asset": paginationResponse})

	case db.AssetTypeERC1155:
		totalAssets, err := db.CountItemsWithFilter(cc.rawDb, "erc_1155_collection_assets", filterConditions)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		assets, _ := db.QueryWithDynamicFilter[db.Erc1155CollectionAsset](cc.rawDb, "erc_1155_collection_assets", limit, offset, filterConditions)

		paginationResponse := db.Pagination[db.Erc1155CollectionAsset]{
			Page:       page,
			Limit:      limit,
			TotalItems: int64(totalAssets),
			TotalPages: int64(totalAssets+(limit)-1) / int64(limit),
			Data:       assets,
		}

		ctx.JSON(http.StatusOK, gin.H{"status": "Successfully retrived id", "type": "ERC721", "asset": paginationResponse})

	case db.AssetTypeERC20:
		totalAssets, err := db.CountItemsWithFilter(cc.rawDb, "erc_20_collection_assets", filterConditions)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		assets, _ := db.QueryWithDynamicFilter[db.Erc20CollectionAsset](cc.rawDb, "erc_20_collection_assets", limit, offset, filterConditions)

		paginationResponse := db.Pagination[db.Erc20CollectionAsset]{
			Page:       page,
			Limit:      limit,
			TotalItems: int64(totalAssets),
			TotalPages: int64(totalAssets+(limit)-1) / int64(limit),
			Data:       assets,
		}

		ctx.JSON(http.StatusOK, gin.H{"status": "Successfully retrived id", "type": "ERC721", "asset": paginationResponse})

	}

}

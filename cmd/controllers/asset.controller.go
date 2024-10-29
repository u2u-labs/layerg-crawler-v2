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
// @Param body body utils.AddNewAssetParamsSwagger true "Asset collection information"
// @Example      { "id": 1, "chain": "U2U", "name": "Nebulas Testnet", "RpcUrl": "sre", "ChainId": 2484, "Explorer": "str", "BlockTime": 500 }
// @Router       /chain/{chain_id}/collection [post]
func (cc *AssetController) AddNewAsset(ctx *gin.Context) {
	// var params *db.AddNewAssetParams
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

	assetParam := utils.ConvertUtilToParams(params)

	// add to db
	if err := cc.db.AddNewAsset(ctx, assetParam); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Asset added", "data": params})
}

// GetAssetByChainId godoc
// @Summary      Get all asset collections for a specific chain
// @Description  Retrieve all asset collections associated with the specified chain ID.
// @Tags         asset
// @Accept       json
// @Produce      json
// @Param chain_id path string true "Chain ID"
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Router       /chain/{chain_id}/collection [get]
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
	paginationResponse := db.Pagination[db.Asset]{
		Page:       page,
		Limit:      limit,
		TotalItems: totalAssets,
		TotalPages: (totalAssets + int64(limit) - 1) / int64(limit), // Calculate total pages
		Data:       assets,
	}

	ctx.JSON(http.StatusOK, gin.H{"data": paginationResponse})
}

// GetAssetCollectionAByChainIdAndContractAddress godoc
// @Summary      Get all asset collection of the chain
// @Description  Get all asset collection of the chain
// @Tags         asset
// @Accept       json
// @Produce      json
// @Param chain_id path string true "Chain ID"
// @Param collection_address path string true "Collection Address"
// @Router       /chain/{chain_id}/collection/{collection_address} [get]
func (cc *AssetController) GetAssetCollectionByChainIdAndContractAddress(ctx *gin.Context) {
	chainIdStr := ctx.Param("chain_id")
	collectionAddress := ctx.Param("collection_address")

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

	response.SuccessReponseData(ctx, http.StatusOK, assetCollection)
}

// GetAssetCollectionAByChainIdAndContractAddress godoc
// @Summary      Get all asset collection of the chain
// @Description  Get all asset collection of the chain
// @Tags         asset
// @Accept       json
// @Produce      json
// @Param chain_id path int true "Chain ID"
// @Param collection_address path string true "Collection Address"
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Param token_id query string false "Token ID"
// @Param owner query string false "Owner Address"
// @Router       /chain/{chain_id}/collection/{collection_address}/asset [get]
func (cc *AssetController) GetAssetByChainIdAndContractAddress(ctx *gin.Context) {
	chainIdStr := ctx.Param("chain_id")
	collectionAddress := ctx.Param("collection_address")

	assetId := chainIdStr + ":" + collectionAddress

	hasFilterParam, tokenId, owner := utils.GetAssetFilterParam(ctx)

	if !hasFilterParam {
		cc.GetAssetsFromCollection(ctx, assetId)
	} else {
		cc.GetAssetsFromCollectionWithFilter(ctx, assetId, tokenId, owner)
	}

}

// GetAssetCollectionAByChainIdAndContractAddress godoc
// @Summary      Get all asset collection of the chain
// @Description  Get all asset collection of the chain
// @Tags         asset
// @Accept       json
// @Produce      json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Param owner query string false "Owner Address"
// @Param asset_type query string false "Asset Type" Enums(ERC721, ERC1155, ERC20)
// @Router       /asset [get]
func (cc *AssetController) GetAssetsByOwner(ctx *gin.Context) {
	owner := ctx.Query("owner")

	page, limit, offset := db.GetLimitAndOffset(ctx)

	switch assetType := ctx.Query("asset_type"); assetType {
	case "ERC721":
		erc721Assets, _ := cc.db.GetPaginated721AssetByOwnerAddress(ctx, db.GetPaginated721AssetByOwnerAddressParams{
			Owner:  owner,
			Limit:  int32(limit),
			Offset: int32(offset),
		})

		totalErc721Assets, err := cc.db.Count721AssetByOwnerAddress(ctx, owner)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		paginationResponse := db.Pagination[db.Erc721CollectionAsset]{
			Page:       page,
			Limit:      limit,
			TotalItems: totalErc721Assets,
			TotalPages: (totalErc721Assets + int64(limit) - 1) / int64(limit),
			Data:       erc721Assets,
		}

		ctx.JSON(http.StatusOK, gin.H{"data": paginationResponse})

	case "ERC1155":
		erc1155Assets, _ := cc.db.GetPaginated1155AssetByOwnerAddress(ctx, db.GetPaginated1155AssetByOwnerAddressParams{
			Owner:  owner,
			Limit:  int32(limit),
			Offset: int32(offset),
		})

		totalErc1155Assets, err := cc.db.Count1155AssetByOwner(ctx, owner)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		paginationResponse := db.Pagination[db.Erc1155CollectionAsset]{
			Page:       page,
			Limit:      limit,
			TotalItems: totalErc1155Assets,
			TotalPages: (totalErc1155Assets + int64(limit) - 1) / int64(limit),
			Data:       erc1155Assets,
		}

		ctx.JSON(http.StatusOK, gin.H{"data": paginationResponse})

	case "ERC20":
		erc20Assets, _ := cc.db.GetPaginated20AssetByOwnerAddress(ctx, db.GetPaginated20AssetByOwnerAddressParams{
			Owner:  owner,
			Limit:  int32(limit),
			Offset: int32(offset),
		})

		totalErc20Assets, err := cc.db.Count20AssetByOwner(ctx, owner)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		paginationResponse := db.Pagination[db.Erc20CollectionAsset]{
			Page:       page,
			Limit:      limit,
			TotalItems: totalErc20Assets,
			TotalPages: (totalErc20Assets + int64(limit) - 1) / int64(limit),
			Data:       erc20Assets,
		}

		ctx.JSON(http.StatusOK, gin.H{"data": paginationResponse})
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

		paginationResponse := db.Pagination[db.Erc721CollectionAsset]{
			Page:       page,
			Limit:      limit,
			TotalItems: totalAssets,
			TotalPages: (totalAssets + int64(limit) - 1) / int64(limit),
			Data:       assets,
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

		paginationResponse := db.Pagination[db.Erc1155CollectionAsset]{
			Page:       page,
			Limit:      limit,
			TotalItems: totalAssets,
			TotalPages: (totalAssets + int64(limit) - 1) / int64(limit),
			Data:       assets,
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

func (cc *AssetController) GetAssetsFromCollectionWithFilter(ctx *gin.Context, assetId string, tokenId string, owner string) {
	assetCollection, err := cc.db.GetAssetById(ctx, assetId)
	filterConditions := make(map[string]string)

	if filterField := ctx.Query("token_id"); filterField != "" {
		filterConditions["token_id"] = filterField
	}

	if filterField := ctx.Query("owner"); filterField != "" {
		filterConditions["owner"] = filterField
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

		paginationResponse := db.Pagination[db.Erc721CollectionAsset]{
			Page:       page,
			Limit:      limit,
			TotalItems: int64(totalAssets),
			TotalPages: int64(totalAssets+(limit)-1) / int64(limit),
			Data:       assets,
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

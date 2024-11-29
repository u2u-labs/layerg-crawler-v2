package controllers

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/unicornultrafoundation/go-u2u/common"

	"github.com/u2u-labs/layerg-crawler/cmd/services"
	"github.com/u2u-labs/layerg-crawler/cmd/utils"
)

type AssetController struct {
	service *services.AssetService
	ctx     context.Context
	rdb     *redis.Client
}

func NewAssetController(service *services.AssetService, ctx context.Context, rdb *redis.Client) *AssetController {
	return &AssetController{service, ctx, rdb}
}

// AddNewAsset godoc
// @Summary      Add a new asset collection to the chain
// @Description  Add a new asset collection to the chain
// @Tags         asset
// @Accept       json
// @Produce      json
// @Param chain_id path string true "Chain Id"
// @Security     ApiKeyAuth
// @Param body body utils.AddNewAssetParamsSwagger true "Asset collection information"
// @Example      { "id": 1, "chain": "U2U", "name": "Nebulas Testnet", "RpcUrl": "sre", "ChainId": 2484, "Explorer": "str", "BlockTime": 500 }
// @Router       /chain/{chain_id}/collection [post]
func (ac *AssetController) AddAssetCollection(ctx *gin.Context) {
	ac.service.AddNewAsset(ctx)
}

// GetAssetCollection godoc
// @Summary      Get all asset in a collection of the chain
// @Description  Retrieve all asset collections associated with the specified chain Id.
// @Tags         asset
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param chain_id path string true "Chain Id"
// @Param collection_address query string false "Collection Address"
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Router       /chain/{chain_id}/collection [get]
func (ac *AssetController) GetAssetCollection(ctx *gin.Context) {

	hasQuery, collectionAddress := utils.GetAssetCollectionFilterParam(ctx)
	if !hasQuery {
		ac.service.GetAssetByChainId(ctx)
	} else {
		ac.service.GetAssetCollectionByChainIdAndContractAddress(ctx, collectionAddress)
	}

}

// GetAssetCollectionAByChainIdAndContractAddress godoc
// @Summary      Get all asset collection of the chain
// @Description  Get all asset collection of the chain
// @Tags         asset
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param chain_id path int true "Chain Id"
// @Param collection_address path string true "Collection Address"
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Param token_id query []string false "Token Ids" collectionFormat(multi)
// @Param owner query string false "Owner Address"
// @Router       /chain/{chain_id}/collection/{collection_address}/assets [get]
func (ac *AssetController) GetAssetByChainIdAndContractAddress(ctx *gin.Context) {
	chainIdStr := ctx.Param("chain_id")
	collectionAddress := ctx.Param("collection_address")
	collectionAddress = common.HexToAddress(collectionAddress).Hex()
	assetId := chainIdStr + ":" + collectionAddress
	hasFilterParam, tokenIds, owner := utils.GetAssetFilterParam(ctx)

	if !hasFilterParam {
		ac.service.GetAssetsFromAssetCollectionId(ctx, assetId)
	} else {
		ac.service.GetAssetsFromCollectionWithFilter(ctx, assetId, tokenIds, owner)
	}

}

// GetNFTAssetCollectionByChainId godoc
// @Summary      Get all NFT asset collection of the chain
// @Description  Get all NFT asset collection of the chain
// @Tags         asset
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param chain_id path int true "Chain Id"
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Router       /chain/{chain_id}/nft-assets [get]
func (ac *AssetController) GetNFTCombinedAsset(ctx *gin.Context) {
	ac.service.GetNFTCombinedAsset(ctx)
}

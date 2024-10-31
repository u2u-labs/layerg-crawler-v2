package utils

import (
	"github.com/gin-gonic/gin"
)

// get asset filter parameters
func GetAssetFilterParam(ctx *gin.Context) (bool, []string, string) {
	// Get the query parameters
	tokenIds := ctx.QueryArray("token_id")
	owner := ctx.Query("owner")

	hasQuery := false
	if len(tokenIds) != 0 || owner != "" {
		hasQuery = true
	}

	return hasQuery, tokenIds, owner
}

// get asset collection filter parameters
func GetAssetCollectionFilterParam(ctx *gin.Context) (bool, string) {
	// Get the query parameters
	collectionAddress := ctx.Query("collection_address")

	hasQuery := false
	if collectionAddress != "" {
		hasQuery = true
	}

	return hasQuery, collectionAddress
}

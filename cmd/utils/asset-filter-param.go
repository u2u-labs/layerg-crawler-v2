package utils

import (
	"github.com/gin-gonic/gin"
)

// get limit and offset from query parameters
func GetAssetFilterParam(ctx *gin.Context) (bool, string, string) {
	// Get the query parameters
	tokenId := ctx.Query("token_id")

	owner := ctx.Query("owner")

	hasQuery := false
	if tokenId != "" || owner != "" {
		hasQuery = true
	}

	return hasQuery, tokenId, owner
}

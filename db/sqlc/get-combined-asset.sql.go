package db

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetCombinedAssetQueryScript
func GetCombinedNFTAssetQueryScript(ctx *gin.Context, limit int, offset int) (string, []interface{}) {
	whereClause, args := GetCombinedAssetWhereClause(ctx)

	// Get the query parameters
	query := fmt.Sprintf(`
	SELECT * FROM (
		SELECT 
			'ERC721' as token_type,
			MIN(chain_id) as chain_id,
			asset_id,
			token_id,
			attributes,
			MIN(created_at) as created_at
			
		FROM erc_721_collection_assets
		%s
		GROUP BY asset_id, token_id, attributes
		UNION ALL
		SELECT 
			'ERC1155' as token_type,
			MIN(chain_id) as chain_id,
			asset_id,
			token_id,
			attributes,
			MIN(created_at) as created_at
			
		FROM erc_1155_collection_assets
		%s
		GROUP BY asset_id, token_id, attributes
	) combined
	ORDER BY combined.created_at ASC
	LIMIT %d OFFSET %d
	
	`, whereClause, whereClause, limit, offset)
	return query, args
}

// GetCombinedAssetQueryScript
func GeCountCombinedNFTAssetQueryScript(ctx *gin.Context) (string, []interface{}) {
	whereClause, args := GetCombinedAssetWhereClause(ctx)

	// Get the query parameters
	query := fmt.Sprintf(`
	SELECT COUNT(*) FROM (
		SELECT 
			'ERC721' as token_type,
			MIN(chain_id) as chain_id,
			asset_id,
			token_id,
			attributes,
			MIN(created_at) as created_at
			
		FROM erc_721_collection_assets
		%s
		GROUP BY asset_id, token_id, attributes
		UNION ALL
		SELECT 
			'ERC1155' as token_type,
			MIN(chain_id) as chain_id,
			asset_id,
			token_id,
			attributes,
			MIN(created_at) as created_at
			
		FROM erc_1155_collection_assets
		%s
		GROUP BY asset_id, token_id, attributes
	) 
	`, whereClause, whereClause)
	return query, args
}

// GetCombinedAssetWhereClause
func GetCombinedAssetWhereClause(ctx *gin.Context) (string, []interface{}) {
	chainIdStr := ctx.Param("chain_id")

	args := []interface{}{}
	whereClause := "WHERE chain_id = $1"

	if chainIdStr != "" {
		chainId, err := strconv.Atoi(chainIdStr)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return "", nil
		}
		args = append(args, chainId)
	}

	created_at_from := ctx.Query("created_at_from")
	created_at_to := ctx.Query("created_at_to")
	if created_at_from != "" {
		whereClause += " AND created_at >= $2"
		args = append(args, created_at_from)
	}

	if created_at_to != "" {
		whereClause += " AND created_at <= $3"
		args = append(args, created_at_to)
	}
	return whereClause, args
}

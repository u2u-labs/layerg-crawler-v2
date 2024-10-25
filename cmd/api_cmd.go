package cmd

import (
	"context"
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/u2u-labs/layerg-crawler/cmd/controllers"
	dbCon "github.com/u2u-labs/layerg-crawler/db/sqlc"
)

func startApi(cmd *cobra.Command, args []string) {

	conn, err := sql.Open(
		viper.GetString("COCKROACH_DB_DRIVER"),
		viper.GetString("COCKROACH_DB_URL"),
	)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	db := dbCon.New(conn)

	if err != nil {
		panic(err)
	}

	serveApi(db, context.Background())
}

func serveApi(db *dbCon.Queries, ctx context.Context) {

	// Create a default Gin router
	router := gin.Default()

	// new Controller
	assetController := controllers.NewAssetController(db, ctx)
	chainController := controllers.NewChainController(db, ctx)

	// Chain routes
	router.POST("/chain", chainController.AddNewChain)
	router.GET("/chain", chainController.GetAllChains)

	// Asset routes
	router.POST("/chain/:chainId/asset", assetController.AddNewAsset)

	router.GET("/chain/:chainId/asset", assetController.GetAssetByChainIdAndContractAddress)

	// Run the server
	router.Run(viper.GetString("API_PORT"))

}

package cmd

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/u2u-labs/layerg-crawler/docs"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/u2u-labs/layerg-crawler/cmd/controllers"
	dbCon "github.com/u2u-labs/layerg-crawler/db/sqlc"
)

// gin-swagger middleware
// swagger embed files

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

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8085

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
// @schemes http https
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

	router.GET("/chain/:chainId/asset", assetController.GetAssetByChainId)

	// Run the server

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(viper.GetString("API_PORT"))
}

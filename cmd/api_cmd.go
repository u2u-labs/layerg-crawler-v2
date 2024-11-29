package cmd

import (
	"context"
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/u2u-labs/layerg-crawler/cmd/controllers"
	middleware "github.com/u2u-labs/layerg-crawler/cmd/middlewares"
	"github.com/u2u-labs/layerg-crawler/cmd/services"
	"github.com/u2u-labs/layerg-crawler/db"
	dbCon "github.com/u2u-labs/layerg-crawler/db/sqlc"
	_ "github.com/u2u-labs/layerg-crawler/docs"
)

func startApi(cmd *cobra.Command, args []string) {
	conn, err := sql.Open(
		viper.GetString("COCKROACH_DB_DRIVER"),
		viper.GetString("COCKROACH_DB_URL"),
	)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	q := dbCon.New(conn)
	if err != nil {
		panic(err)
	}

	rdb, err := db.NewRedisClient(&db.RedisConfig{
		Url:      viper.GetString("REDIS_DB_URL"),
		Db:       viper.GetInt("REDIS_DB"),
		Password: viper.GetString("REDIS_DB_PASSWORD"),
	})
	if err != nil {
		panic(err)
	}

	serveApi(q, rdb, conn, context.Background())
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

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
// @schemes http https
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-KEY
// @Security ApiKeyAuth
func serveApi(db *dbCon.Queries, rdb *redis.Client, rawDb *sql.DB, ctx context.Context) {
	// Create a default Gin router
	gin.SetMode(viper.GetString("GIN_MODE"))
	router := gin.Default()

	// new Service
	chainService := services.NewChainService(db, rawDb, ctx, rdb)
	assetService := services.NewAssetService(db, rawDb, ctx, rdb)

	// new Controller
	chainController := controllers.NewChainController(chainService, ctx, rdb)
	assetController := controllers.NewAssetController(assetService, ctx, rdb)
	historyController := controllers.NewHistoryController(db, rawDb, ctx, rdb)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Apply the basic authentication middleware
	router.Use(middleware.ApiKeyAuth(db))

	// Chain routes
	router.POST("/chain", chainController.AddNewChain)
	router.GET("/chain", chainController.GetAllChains)

	// Asset routes
	router.POST("/chain/:chain_id/collection", assetController.AddAssetCollection)
	router.GET("/chain/:chain_id/collection", assetController.GetAssetCollection)
	router.GET("/chain/:chain_id/collection/:collection_address/assets", assetController.GetAssetByChainIdAndContractAddress)
	router.GET("/chain/:chain_id/nft-assets", assetController.GetNFTCombinedAsset)

	// History routes
	router.GET("/history", historyController.GetHistory)

	// Run the server

	router.Run(viper.GetString("API_PORT"))
}

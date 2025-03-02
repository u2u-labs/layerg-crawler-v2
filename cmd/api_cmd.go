package cmd

import (
	_ "github.com/lib/pq"

	_ "github.com/u2u-labs/layerg-crawler/docs"
)

// func startApi(cmd *cobra.Command, args []string) {
// 	var (
// 		logger = &zap.Logger{}
// 	)
// 	conn, err := sql.Open(
// 		viper.GetString("COCKROACH_DB_DRIVER"),
// 		viper.GetString("COCKROACH_DB_URL"),
// 	)
// 	sugar := logger.Sugar()

// 	if err != nil {
// 		sugar.Errorw("Failed to connect to database", "err", err)
// 	}

// 	q := dbCon.New(conn)

// 	rdb, err := db.NewRedisClient(&db.RedisConfig{
// 		Url:      viper.GetString("REDIS_DB_URL"),
// 		Db:       viper.GetInt("REDIS_DB"),
// 		Password: viper.GetString("REDIS_DB_PASSWORD"),
// 	})
// 	if err != nil {
// 		sugar.Errorw("Failed to connect to redis", "err", err)
// 	}

// 	serveApi(q, rdb, conn, context.Background())
// }

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
// func serveApi(db *dbCon.Queries, rdb *redis.Client, rawDb *sql.DB, ctx context.Context) {
// 	// Create a default Gin router
// 	gin.SetMode(viper.GetString("GIN_MODE"))
// 	router := gin.Default()

// 	// new Service
// 	chainService := services.NewChainService(db, rawDb, ctx, rdb)
// 	assetService := services.NewAssetService(db, rawDb, ctx, rdb)

// 	// new Controller
// 	chainController := controllers.NewChainController(chainService, ctx, rdb)
// 	assetController := controllers.NewAssetController(assetService, ctx, rdb)
// 	historyController := controllers.NewHistoryController(db, rawDb, ctx, rdb)
// 	backfillController := controllers.NewBackFillController(db, rawDb, ctx, rdb)

// 	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

// 	// Apply the basic authentication middleware
// 	router.Use(middleware.ApiKeyAuth(db))

// 	// Chain routes
// 	router.POST("/chain", chainController.AddNewChain)
// 	router.GET("/chain", chainController.GetAllChains)

// 	// Asset routes
// 	router.POST("/chain/:chain_id/collection", assetController.AddAssetCollection)
// 	router.GET("/chain/:chain_id/collection", assetController.GetAssetCollection)
// 	router.GET("/chain/:chain_id/collection/:collection_address/assets", assetController.GetAssetByChainIdAndContractAddress)
// 	router.GET("/chain/:chain_id/collection/:collection_address/:token_id", assetController.GetAssetByChainIdAndContractAddressDetail)
// 	router.GET("/chain/:chain_id/nft-assets", assetController.GetNFTCombinedAsset)

// 	// Backfill routes
// 	router.POST("/backfill", backfillController.AddBackFillTracker)
// 	// History routes``
// 	router.GET("/history", historyController.GetHistory)

// 	// Run the server

// 	router.Run(viper.GetString("API_PORT"))
// }

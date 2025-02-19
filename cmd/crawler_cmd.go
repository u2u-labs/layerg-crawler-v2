package cmd

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/hibiken/asynq"
	"github.com/unicornultrafoundation/go-u2u/ethclient"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/google/uuid"
	"github.com/u2u-labs/layerg-crawler/config"
	"github.com/u2u-labs/layerg-crawler/db"
	"github.com/u2u-labs/layerg-crawler/db/graphqldb"
	dbCon "github.com/u2u-labs/layerg-crawler/db/sqlc"
	"github.com/u2u-labs/layerg-crawler/generated/router"
)

var contractType = make(map[int32]map[string]dbCon.Asset)

func startCrawler(cmd *cobra.Command, args []string) {
	var (
		ctx    = context.Background()
		logger = &zap.Logger{}
	)
	if viper.GetString("LOG_LEVEL") == "PROD" {
		logger, _ = zap.NewProduction()
	} else {
		logger, _ = zap.NewDevelopment()
	}
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	conn, err := sql.Open(
		viper.GetString("COCKROACH_DB_DRIVER"),
		viper.GetString("COCKROACH_DB_URL"),
	)
	if err != nil {
		sugar.Errorw("Could not connect to database", "err", err)
		return
	}

	sqlDb := dbCon.New(conn)

	// Load crawler config
	crawlerConfig, err := loadCrawlerConfig()
	if err != nil {
		sugar.Errorw("Failed to load crawler config", "err", err)
		return
	}

	// Initialize chain from subgraph config
	chainID, _ := strconv.ParseInt(crawlerConfig.Network.ChainId, 10, 32)
	chainIDInt32 := int32(chainID)

	// Find the lowest startBlock from all datasources
	var lowestStartBlock int64
	if len(crawlerConfig.DataSources) > 0 {
		lowestStartBlock = crawlerConfig.DataSources[0].StartBlock
		for _, ds := range crawlerConfig.DataSources {
			if ds.StartBlock < lowestStartBlock {
				lowestStartBlock = ds.StartBlock
			}
		}
	}

	// Create or update chain information
	_, err = sqlDb.CreateChain(ctx, dbCon.CreateChainParams{
		ID:          chainIDInt32,
		Chain:       "ethereum",
		Name:        crawlerConfig.Network.Name,
		RpcUrl:      crawlerConfig.Network.Endpoint[0],
		ChainID:     chainID,
		Explorer:    "",
		LatestBlock: lowestStartBlock,
		BlockTime:   500,
	})
	if err != nil {
		sugar.Warnw("Failed to create chain", "err", err)
		// return
	}

	// Initialize assets (contracts) from subgraph config
	for _, ds := range crawlerConfig.DataSources {
		startBlock := sql.NullInt64{
			Int64: ds.StartBlock,
			Valid: true,
		}

		_, err = sqlDb.CreateAsset(ctx, dbCon.CreateAssetParams{
			ID:              uuid.New().String(),
			ChainID:         chainIDInt32,
			ContractAddress: strings.ToLower(ds.Options.Address),
			InitialBlock:    startBlock,
		})
		if err != nil {
			sugar.Warnw("Failed to create asset",
				"err", err,
				"address", ds.Options.Address,
				"startBlock", ds.StartBlock,
			)
			// return
		}
		sugar.Infow("Initialized contract",
			"address", ds.Options.Address,
			"startBlock", ds.StartBlock,
		)
	}

	router := router.NewEventRouter(sqlDb, graphqldb.New(conn), sugar, chainIDInt32)

	// Initialize Redis and continue with existing code...
	rdb, err := db.NewRedisClient(&db.RedisConfig{
		Url:      viper.GetString("REDIS_DB_URL"),
		Db:       viper.GetInt("REDIS_DB"),
		Password: viper.GetString("REDIS_DB_PASSWORD"),
	})
	if err != nil {
		sugar.Errorw("Failed to connect to redis", "err", err)
	}

	queueClient := asynq.NewClient(asynq.RedisClientOpt{Addr: viper.GetString("REDIS_DB_URL")})
	defer queueClient.Close()

	err = crawlSupportedChains(ctx, sugar, sqlDb, rdb, router)
	if err != nil {
		sugar.Errorw("Error init supported chains", "err", err)
		return
	}

	// TODO: for each endpoint, start a separate worker
	rpcURL := crawlerConfig.Network.Endpoint[0]

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		sugar.Errorw("Failed to connect to ethereum client", "err", err, "rpc", rpcURL)
		return
	}
	defer client.Close()

	timer := time.NewTimer(config.RetriveAddedChainsAndAssetsInterval)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			// Process new chains with the event registry
			ProcessNewChains(ctx, sugar, rdb, sqlDb, router)
			// Process new assets
			ProcessNewChainAssets(ctx, sugar, rdb)
			timer.Reset(config.RetriveAddedChainsAndAssetsInterval)
		case <-ctx.Done():
			return
		}
	}
}

func ProcessNewChains(ctx context.Context, sugar *zap.SugaredLogger, rdb *redis.Client, q *dbCon.Queries, registry *router.EventRouter) {
	chains, err := db.GetCachedPendingChain(ctx, rdb)
	if err != nil {
		sugar.Errorw("ProcessNewChains failed to get cached pending chains", "err", err)
		return
	}
	if err = db.DeletePendingChainsInCache(ctx, rdb); err != nil {
		sugar.Errorw("ProcessNewChains failed to delete cached pending chains", "err", err)
		return
	}
	for _, c := range chains {
		client, err := initChainClient(&c)

		if err != nil {
			sugar.Errorw("ProcessNewChains failed to init chain client", "err", err, "chain", c)
			return
		}

		go StartChainCrawler(ctx, sugar, client, q, &c, registry)
		sugar.Infow("Initiated new chain, start crawling", "chain", c)
	}
}

func ProcessNewChainAssets(ctx context.Context, sugar *zap.SugaredLogger, rdb *redis.Client) {
	assets, err := db.GetCachedPendingAsset(ctx, rdb)
	if err != nil {
		sugar.Errorw("ProcessNewChainAssets failed to get cached pending assets", "err", err)
		return
	}
	if err = db.DeletePendingAssetsInCache(ctx, rdb); err != nil {
		sugar.Errorw("ProcessNewChainAssets failed to delete cached pending assets", "err", err)
		return
	}
	for _, a := range assets {
		contractType[a.ChainID][a.ContractAddress] = a
		sugar.Infow("Initiated new assets, start crawling",
			"chain", a.ChainID,
			"address", a.ContractAddress,
		)
	}
}

func crawlSupportedChains(ctx context.Context, sugar *zap.SugaredLogger, q *dbCon.Queries, rdb *redis.Client, registry *router.EventRouter) error {
	// Query, flush cache and connect all supported chains
	chains, err := q.GetAllChain(ctx)
	if err != nil {
		return err
	}
	for _, c := range chains {
		contractType[c.ID] = make(map[string]dbCon.Asset)
		if err = db.DeleteChainInCache(ctx, rdb, c.ID); err != nil {
			return err
		}
		if err = db.DeleteChainAssetsInCache(ctx, rdb, c.ID); err != nil {
			return err
		}

		// Query all assets of one chain
		var (
			assets []dbCon.Asset
			limit  int32 = 10
			offset int32 = 0
		)
		for {
			a, err := q.GetPaginatedAssetsByChainId(ctx, dbCon.GetPaginatedAssetsByChainIdParams{
				ChainID: c.ID,
				Limit:   limit,
				Offset:  offset,
			})
			if err != nil {
				return err
			}
			assets = append(assets, a...)
			offset = offset + limit
			if len(a) < int(limit) {
				break
			}
		}

		if err = db.SetChainToCache(ctx, rdb, &c); err != nil {
			return err
		}
		if err = db.SetAssetsToCache(ctx, rdb, assets); err != nil {
			return err
		}
		for _, a := range assets {
			contractType[a.ChainID][a.ContractAddress] = a
		}
		client, err := initChainClient(&c)
		if err != nil {
			return err
		}
		go StartChainCrawler(ctx, sugar, client, q, &c, registry)
	}
	return nil
}

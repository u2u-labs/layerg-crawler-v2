package cmd

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/hibiken/asynq"
	"github.com/u2u-labs/layerg-crawler/config"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/unicornultrafoundation/go-u2u/accounts/abi"
	"go.uber.org/zap"

	"github.com/u2u-labs/layerg-crawler/cmd/utils"
	"github.com/u2u-labs/layerg-crawler/db"
	dbCon "github.com/u2u-labs/layerg-crawler/db/sqlc"
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
	}

	sqlDb := dbCon.New(conn)

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

	if utils.ERC20ABI, err = abi.JSON(strings.NewReader(utils.ERC20ABIStr)); err != nil {
		sugar.Errorw("Failed to parse ERC20 ABI", "err", err)
	}
	if utils.ERC721ABI, err = abi.JSON(strings.NewReader(utils.ERC721ABIStr)); err != nil {
		sugar.Errorw("Failed to parse ERC721 ABI", "err", err)
	}
	if utils.ERC1155ABI, err = abi.JSON(strings.NewReader(utils.ERC1155ABIStr)); err != nil {
		sugar.Errorw("Failed to parse ERC1155 ABI", "err", err)
	}

	err = crawlSupportedChains(ctx, sugar, sqlDb, rdb)
	if err != nil {
		sugar.Errorw("Error init supported chains", "err", err)
		return
	}

	timer := time.NewTimer(config.RetriveAddedChainsAndAssetsInterval)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			// Process new chains
			ProcessNewChains(ctx, sugar, rdb, sqlDb)
			// Process new assets
			ProcessNewChainAssets(ctx, sugar, rdb)
			// Process backfill collection
			ProcessCrawlingBackfillCollection(ctx, sugar, sqlDb, rdb, queueClient)

			timer.Reset(config.RetriveAddedChainsAndAssetsInterval)
		}
	}
}

func ProcessNewChains(ctx context.Context, sugar *zap.SugaredLogger, rdb *redis.Client, q *dbCon.Queries) {
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

		go StartChainCrawler(ctx, sugar, client, q, &c, rdb)
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
		contractType[a.ChainID][a.CollectionAddress] = a
		sugar.Infow("Initiated new assets, start crawling",
			"chain", a.ChainID,
			"address", a.CollectionAddress,
			"type", a.Type,
		)
	}
}

func crawlSupportedChains(ctx context.Context, sugar *zap.SugaredLogger, q *dbCon.Queries, rdb *redis.Client) error {
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
			contractType[a.ChainID][a.CollectionAddress] = a
		}
		client, err := initChainClient(&c)
		if err != nil {
			return err
		}
		go StartChainCrawler(ctx, sugar, client, q, &c, rdb)

	}
	return nil

}

func ProcessCrawlingBackfillCollection(ctx context.Context, sugar *zap.SugaredLogger, q *dbCon.Queries, rdb *redis.Client, queueClient *asynq.Client) error {
	// Get all Backfill Collection with status CRAWLING
	crawlingBackfill, err := q.GetCrawlingBackfillCrawler(ctx)

	if err != nil {
		return err
	}

	for _, c := range crawlingBackfill {
		chain, err := q.GetChainById(ctx, c.ChainID)

		client, err := initChainClient(&chain)
		if err != nil {
			return err
		}
		go AddBackfillCrawlerTask(ctx, sugar, client, q, &chain, &c, queueClient)

	}
	return nil
}

package cmd

import (
	"context"
	"database/sql"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/u2u-labs/layerg-crawler/db"
	dbCon "github.com/u2u-labs/layerg-crawler/db/sqlc"
)

func startCrawler(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	conn, err := sql.Open(
		viper.GetString("COCKROACH_DB_DRIVER"),
		viper.GetString("COCKROACH_DB_URL"),
	)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	sqlDb := dbCon.New(conn)

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

	err = crawlSupportedChains(ctx, sugar, sqlDb, rdb)
	if err != nil {
		sugar.Errorw("Error init supported chains", "err", err)
		return
	}
	select {}
}

func crawlSupportedChains(ctx context.Context, sugar *zap.SugaredLogger, q *dbCon.Queries, rdb *redis.Client) error {
	// Query, cache and connect all supported chains
	chains, err := q.GetAllChain(ctx)
	if err != nil {
		return err
	}
	for _, c := range chains {
		_, err := q.GetAssetByChainId(ctx, c.ID)
		if err != nil {
			return err
		}

		if err := db.SetChainToCache(rdb, &c); err != nil {
			return err
		}
		client, err := initChainClient(&c)
		if err != nil {
			return err
		}
		go StartChainCrawler(sugar, client, &c)
	}

	return nil
}

package cmd

import (
	"context"
	"database/sql"
	"github.com/u2u-labs/layerg-crawler/cmd/utils"
	"log"
	"strings"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/unicornultrafoundation/go-u2u/accounts/abi"
	"go.uber.org/zap"

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

	if utils.ERC20ABI, err = abi.JSON(strings.NewReader(utils.ERC20ABIStr)); err != nil {
		panic(err)
	}
	if utils.ERC721ABI, err = abi.JSON(strings.NewReader(utils.ERC721ABIStr)); err != nil {
		panic(err)
	}
	if utils.ERC1155ABI, err = abi.JSON(strings.NewReader(utils.ERC1155ABIStr)); err != nil {
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
		go StartChainCrawler(ctx, sugar, client, q, &c)
	}
	return nil

}

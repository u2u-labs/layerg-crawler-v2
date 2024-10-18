package cmd

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/u2u-labs/layerg-crawler/db"
	"github.com/u2u-labs/layerg-crawler/types"
)

func startCrawler(cmd *cobra.Command, args []string) {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	gdb, err := db.NewCockroachDbClient(&db.DbConfig{
		Url:  viper.GetString("COCKROACH_DB_URL"),
		Name: viper.GetString("COCKROACH_DB_NAME"),
	})
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

	err = initSupportedChains(sugar, gdb, rdb)
	if err != nil {
		sugar.Errorw("Error init supported chains", "err", err)
		return
	}
	select {}
}

func initSupportedChains(sugar *zap.SugaredLogger, gdb *gorm.DB, rdb *redis.Client) error {
	if !gdb.Migrator().HasTable(&types.Chain{}) {
		if err := gdb.AutoMigrate(&types.Chain{}); err != nil {
			return err
		}
		if err := db.InsertSupportedChains(gdb); err != nil {
			return err
		}
	}
	var (
		chains []*types.Chain
	)
	// Query, cache and connect all supported chains
	gdb.Find(&chains)
	for _, chain := range chains {
		if err := db.SetChainToCache(rdb, chain); err != nil {
			return err
		}
		client, err := initChainClient(chain)
		if err != nil {
			return err
		}
		go StartChainCrawler(sugar, client, chain)
	}
	return nil
}

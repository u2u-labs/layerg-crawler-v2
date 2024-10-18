package cmd

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/gorm"

	"github.com/u2u-labs/layerg-crawler/db"
	"github.com/u2u-labs/layerg-crawler/types"
)

func startCrawler(cmd *cobra.Command, args []string) {
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
	err = initSupportedChains(gdb, rdb)
	if err != nil {
		panic(err)
	}
}

func initSupportedChains(gdb *gorm.DB, rdb *redis.Client) error {
	if err := gdb.AutoMigrate(&types.Network{}); err != nil {
		return err
	}
	if err := db.InsertSupportedChains(gdb); err != nil {
		return err
	}
	var (
		chains []*types.Network
	)
	// Query, cache and connect all supported chains
	gdb.Find(&chains)
	for _, chain := range chains {
		if err := db.SetChain(rdb, chain); err != nil {
			return err
		}
		c, err := initChainClient(chain)
		if err != nil {
			return err
		}
		StartChainCrawler(c, chain)
	}
	return nil
}

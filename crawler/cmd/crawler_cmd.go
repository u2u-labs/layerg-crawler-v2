package cmd

import (
	"context"
	"strconv"
	"time"

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
	initSupportedChains(gdb, rdb)
}

func initSupportedChains(gdb *gorm.DB, rdb *redis.Client) {
	gdb.AutoMigrate(&types.Network{})
	db.InsertSupportedChains(gdb)
	var (
		ctx    = context.Background()
		chains []*types.Network
	)
	gdb.Find(chains)
	for _, chain := range chains {
		rdb.Set(ctx, strconv.Itoa(chain.Id), chain.String(), 2*time.Second)
	}
}

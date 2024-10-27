package db

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/redis/go-redis/v9"

	db "github.com/u2u-labs/layerg-crawler/db/sqlc"
)

type RedisConfig struct {
	Url      string
	Db       int
	Password string
}

func NewRedisClient(cfg *RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Url,
		Password: cfg.Password,
		DB:       cfg.Db,
	})
	return rdb, nil
}

var ctx = context.Background()

func ChainCacheKey(chainId int32) string {
	return "chain" + strconv.Itoa(int(chainId))
}

func GetCachedChain(rdb *redis.Client, chainId int32) (*db.Chain, error) {
	res := rdb.Get(ctx, ChainCacheKey(chainId))
	if res.Err() != nil {
		return nil, res.Err()
	}
	var chain *db.Chain
	err := json.Unmarshal([]byte(res.Val()), &chain)
	return chain, err
}

func SetChainToCache(rdb *redis.Client, chain *db.Chain) error {
	jsonChain, err := json.Marshal(chain)
	if err != nil {
		return err
	}
	return rdb.Set(ctx, ChainCacheKey(chain.ID), string(jsonChain), 0).Err()
}

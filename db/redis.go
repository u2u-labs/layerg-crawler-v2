package db

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/redis/go-redis/v9"
	
	"github.com/u2u-labs/layerg-crawler/types"
)

type RedisConfig struct {
	Url string
	Db  int

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

func ChainKey(chainId int) string {
	return "chain" + strconv.Itoa(chainId)
}

func GetChain(rdb *redis.Client, chainId int) (*types.Network, error) {
	res := rdb.Get(ctx, ChainKey(chainId))
	if res.Err() != nil {
		return nil, res.Err()
	}
	var chain *types.Network
	err := json.Unmarshal([]byte(res.Val()), &chain)
	return chain, err
}

func SetChain(rdb *redis.Client, chain *types.Network) error {
	jsonChain, err := json.Marshal(chain)
	if err != nil {
		return err
	}
	return rdb.Set(ctx, ChainKey(chain.Id), string(jsonChain), 0).Err()
}

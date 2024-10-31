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

func ChainCacheKey(chainId int32) string {
	return "chain:" + strconv.Itoa(int(chainId))
}

func AssetCacheKey(chainId int32) string {
	return "assets:" + strconv.Itoa(int(chainId))
}

func GetCachedChain(ctx context.Context, rdb *redis.Client, chainId int32) (*db.Chain, error) {
	res := rdb.Get(ctx, ChainCacheKey(chainId))
	if res.Err() != nil {
		return nil, res.Err()
	}
	var chain *db.Chain
	err := json.Unmarshal([]byte(res.Val()), &chain)
	return chain, err
}

func SetChainToCache(ctx context.Context, rdb *redis.Client, chain *db.Chain) error {
	jsonChain, err := json.Marshal(chain)
	if err != nil {
		return err
	}
	return rdb.Set(ctx, ChainCacheKey(chain.ID), string(jsonChain), 0).Err()
}

func DeleteChainInCache(ctx context.Context, rdb *redis.Client, chainId int32) error {
	return rdb.Del(ctx, ChainCacheKey(chainId)).Err()
}

func GetCachedAssets(ctx context.Context, rdb *redis.Client, chainId int32) ([]db.Asset, error) {
	assetsStr, err := rdb.LRange(ctx, AssetCacheKey(chainId), 0, -1).Result()
	if err != nil {
		return nil, err
	}
	var assets []db.Asset
	for _, a := range assetsStr {
		var asset db.Asset
		err = json.Unmarshal([]byte(a), &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, asset)
	}
	return assets, err
}

func SetAssetsToCache(ctx context.Context, rdb *redis.Client, assets []db.Asset) error {
	for _, a := range assets {
		jsonAsset, err := json.Marshal(a)
		if err != nil {
			return err
		}
		err = rdb.LPush(ctx, AssetCacheKey(a.ChainID), string(jsonAsset), 0).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteChainAssetsInCache(ctx context.Context, rdb *redis.Client, chainId int32) error {
	return rdb.Del(ctx, AssetCacheKey(chainId)).Err()
}

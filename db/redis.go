package db

import (
	"github.com/redis/go-redis/v9"
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

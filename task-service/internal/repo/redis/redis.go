package redis

import (
	"context"
	"fmt"
	"task-service/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisDB struct {
	rdb *redis.Client
}

func NewClient(ctx context.Context, cfg *config.Config) (*RedisDB, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:        cfg.Redis.Host + ":" + fmt.Sprint(cfg.Redis.Port),
		Password:    cfg.Redis.Password,
		DB:          cfg.Redis.DB,
		Username:    cfg.Redis.User,
		MaxRetries:  cfg.Redis.MaxRetries,
		DialTimeout: cfg.Redis.DialTimeout,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("Failed to connect to redis: %s", err.Error())
	}
	return &RedisDB{rdb: rdb}, nil
}

func (r *RedisDB) Get(ctx context.Context, uuid string) (string, error) {
	cached, err := r.rdb.Get(ctx, uuid).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return cached, nil
}

func (r *RedisDB) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return r.rdb.Set(ctx, key, value, ttl).Err()
}

func (r *RedisDB) Delete(ctx context.Context, key string) error {
	return r.rdb.Del(ctx, key).Err()
}

package cache

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v9"

	"github.com/VladPetriv/tg_scanner/config"
)

type redisStore struct {
	cfg    *config.Config
	client *redis.Client
}

func NewRedis(cfg *config.Config) Store {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       0,
	})

	return &redisStore{
		cfg:    cfg,
		client: client,
	}
}

func (r redisStore) Get(ctx context.Context, key string) (string, error) {
	value, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}

		return "", fmt.Errorf("get value by key from redis: %w", err)
	}

	return value, nil
}

func (r redisStore) Set(ctx context.Context, key string, value bool) error {
	err := r.client.Set(ctx, key, value, 0)
	if err.Err() != nil {
		return fmt.Errorf("set value by key to redis: %w", err.Err())
	}

	return nil
}

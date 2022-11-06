package redis

import (
	"context"
	"fmt"
	"strings"

	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/pkg/config"
	"github.com/go-redis/redis/v9"
)

type redisStore struct {
	cfg    *config.Config
	client *redis.Client
}

func New(cfg *config.Config) *redisStore {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       0,
	})

	return &redisStore{cfg: cfg, client: client}
}

func (r redisStore) GenerateKey(value interface{}) string {
	var key string

	switch data := value.(type) {
	case model.TgMessage:
		key = fmt.Sprintf(
			"[%d%d%d-%s]",
			data.ID,
			data.FromID.ID,
			data.PeerID.ID,
			strings.ReplaceAll(data.Message, " ", ""),
		)
	case model.TgGroup:
		key = fmt.Sprintf("[%d%s]", data.ID, data.Username)
	}

	return key
}

func (r redisStore) Get(ctx context.Context, key string) (string, error) {
	value, err := r.client.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return "", fmt.Errorf("get data from redis error: %w", err)
	}

	return value, nil
}

func (r redisStore) Set(ctx context.Context, key string, value bool) error {
	err := r.client.Set(ctx, key, value, 0)
	if err.Err() != nil {
		return fmt.Errorf("set data to redis error: %w", err.Err())
	}

	return nil
}

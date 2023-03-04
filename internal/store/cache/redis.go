package cache

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-redis/redis/v9"

	"github.com/VladPetriv/tg_scanner/config"
	"github.com/VladPetriv/tg_scanner/internal/model"
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

	return &redisStore{cfg: cfg, client: client}
}

func (r redisStore) Get(ctx context.Context, data interface{}) (string, error) {
	value, err := r.client.Get(ctx, generateKey(data)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", fmt.Errorf("get data from redis error: %w", err)
	}

	return value, nil
}

func (r redisStore) Set(ctx context.Context, data interface{}, value bool) error {
	err := r.client.Set(ctx, generateKey(data), value, 0)
	if err.Err() != nil {
		return fmt.Errorf("set data to redis error: %w", err.Err())
	}

	return nil
}

func generateKey(value interface{}) string {
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

package redis

import (
	"context"
	"fmt"
	"strings"

	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/pkg/config"
	"github.com/go-redis/redis/v9"
)

type RedisDB struct {
	cfg    *config.Config
	client *redis.Client
}

func NewRedisDB(config *config.Config) *RedisDB {
	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPassword,
		DB:       0,
	})

	return &RedisDB{cfg: config, client: client}
}

func GenerateMessageKey(message model.TgMessage) string {
	key := fmt.Sprintf(
		"[%d%d%d-%s]",
		message.ID,
		message.FromID.ID,
		message.PeerID.ID,
		strings.ReplaceAll(message.Message, " ", ""),
	)

	return key
}

func (r RedisDB) GetMessageFromRedis(ctx context.Context, key string) (string, error) {
	value, err := r.client.Get(ctx, key).Result()
	if err != redis.Nil {
		return "", fmt.Errorf("failed to get message from redis: %w", err)
	}

	return value, nil
}

func (r RedisDB) SetMessageToRedis(ctx context.Context, key string, value bool) error {
	err := r.client.Set(ctx, key, value, 0)
	if err.Err() != nil {
		return fmt.Errorf("failed to set message to redis: %w", err.Err())
	}

	return nil
}

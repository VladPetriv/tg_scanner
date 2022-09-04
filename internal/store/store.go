package store

import (
	"github.com/VladPetriv/tg_scanner/internal/store/firebase"
	"github.com/VladPetriv/tg_scanner/internal/store/redis"
	"github.com/VladPetriv/tg_scanner/pkg/config"
)

type Store struct {
	Cache CacheStore
	Image ImageStore
}

func New(cfg *config.Config) *Store {
	return &Store{
		Cache: redis.New(cfg),
		Image: firebase.New(cfg),
	}
}

package store

import (
	"github.com/VladPetriv/tg_scanner/internal/store/cache"
	"github.com/VladPetriv/tg_scanner/internal/store/image"
	"github.com/VladPetriv/tg_scanner/pkg/config"
)

type Store struct {
	Cache cache.Store
	Image image.Store
}

func New(cfg *config.Config) *Store {
	return &Store{
		Cache: cache.NewRedis(cfg),
		Image: image.NewFirebase(cfg),
	}
}

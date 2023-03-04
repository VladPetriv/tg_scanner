package store

import (
	"github.com/VladPetriv/tg_scanner/config"
	"github.com/VladPetriv/tg_scanner/internal/store/cache"
	"github.com/VladPetriv/tg_scanner/internal/store/image"
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

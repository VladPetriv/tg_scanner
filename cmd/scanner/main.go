package main

import (
	"sync"

	"github.com/VladPetriv/tg_scanner/internal/client"
	"github.com/VladPetriv/tg_scanner/internal/file"
	"github.com/VladPetriv/tg_scanner/internal/store/redis"
	"github.com/VladPetriv/tg_scanner/pkg/config"
	"github.com/VladPetriv/tg_scanner/pkg/logger"
)

func main() {
	err := file.CreateDirs()
	if err != nil {
		panic(err)
	}

	log := logger.Get()

	cfg, err := config.Get()
	if err != nil {
		log.Panic(err)
	}

	var waitGroup sync.WaitGroup

	redisDB := redis.NewRedisDB(cfg)

	client.Run(redisDB, &waitGroup, cfg, log)
}

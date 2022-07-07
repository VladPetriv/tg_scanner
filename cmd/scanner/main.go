package main

import (
	"sync"

	_ "github.com/lib/pq"

	"github.com/VladPetriv/tg_scanner/config"
	"github.com/VladPetriv/tg_scanner/internal/client"
	"github.com/VladPetriv/tg_scanner/internal/file"
	"github.com/VladPetriv/tg_scanner/internal/service"
	"github.com/VladPetriv/tg_scanner/internal/store"
	"github.com/VladPetriv/tg_scanner/logger"
)

func main() {
	// Create needed dirs
	err := file.CreateDirs()
	if err != nil {
		panic(err)
	}

	// Initialize logger
	log := logger.Get()

	// Initialize config
	cfg, err := config.Get()
	if err != nil {
		log.Panic(err)
	}
	var waitGroup sync.WaitGroup

	// Initialize store
	store, err := store.New(*cfg, log)
	if err != nil {
		log.Error(err)
	}

	// Initialize service manager
	serviceManager, err := service.NewManager(store)
	if err != nil {
		log.Error(err)
	}

	client.Run(serviceManager, &waitGroup, cfg, log)
}

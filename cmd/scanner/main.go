package main

import (
	"log"

	"github.com/VladPetriv/tg_scanner/config"
	"github.com/VladPetriv/tg_scanner/internal/app"
	"github.com/VladPetriv/tg_scanner/internal/controller"
	"github.com/VladPetriv/tg_scanner/internal/store"
	"github.com/VladPetriv/tg_scanner/pkg/file"
	"github.com/VladPetriv/tg_scanner/pkg/logger"
)

func main() {
	err := file.InitDirectories()
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := config.Get()
	if err != nil {
		log.Fatal(err)
	}

	logger := logger.Get(cfg)

	store := store.New(cfg)

	queue := controller.New(cfg)

	app.Run(store, queue, cfg, logger)
}

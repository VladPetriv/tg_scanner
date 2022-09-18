package main

import (
	"github.com/VladPetriv/tg_scanner/internal/app"
	"github.com/VladPetriv/tg_scanner/internal/store"
	"github.com/VladPetriv/tg_scanner/pkg/config"
	"github.com/VladPetriv/tg_scanner/pkg/file"
	"github.com/VladPetriv/tg_scanner/pkg/logger"
)

func main() {
	err := file.InitDirectories()
	if err != nil {
		panic(err)
	}

	cfg, err := config.Get()
	if err != nil {
		panic(err)
	}

	log := logger.Get(cfg)

	store := store.New(cfg)

	app.Run(store, cfg, log)
}

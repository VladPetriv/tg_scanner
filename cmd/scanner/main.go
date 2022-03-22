package main

import (
	"context"
	"fmt"
	"github.com/VladPetriv/tg_scanner/config"
	"github.com/VladPetriv/tg_scanner/internal/auth"
	"github.com/VladPetriv/tg_scanner/internal/channel"
	"github.com/VladPetriv/tg_scanner/internal/client"
	"github.com/VladPetriv/tg_scanner/internal/file"
	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/internal/service"
	"github.com/VladPetriv/tg_scanner/internal/store"
	"github.com/VladPetriv/tg_scanner/logger"
	"github.com/gotd/td/telegram"
	"sync"
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

	store, err := store.New(*cfg, log)
	if err != nil {
		log.Error(err)

	}

	serviceManager, err := service.NewManager(store)
	if err != nil {
		log.Error(err)
	}

	// Create new client
	tgClient, err := telegram.ClientFromEnvironment(telegram.Options{}) // nolint
	if err != nil {
		log.Errorf("ERROR_WHILE_CREATING_CLIENT:%s", err)
	}

	// Create API
	api := tgClient.API()

	if err := tgClient.Run(context.Background(), func(ctx context.Context) error {
		// Authorization to telegram
		user, err := auth.Login(ctx, tgClient, cfg)
		if err != nil {
			return fmt.Errorf("AUTH_ERROR:%w", err)
		}

		waitGroup.Add(2) // nolint
		// Get user data
		u, _ := user.GetUser().AsNotEmpty()

		// Getting incoming messages
		go client.GetNewMessage(ctx, u, api, &waitGroup, log)

		// Getting all groups
		groups, err := channel.GetAllGroups(ctx, api)
		if err != nil {
			return fmt.Errorf("GROUPS_ERROR:%w", err)
		}

		// Create files for groups
		file.CreateFilesForGroups(groups)

		// Getting group history
		for _, group := range groups {
			err := serviceManager.Channel.CreateChannel(&model.Channel{Name: group.Title})
			if err != nil {
				log.Error(err)
			}
			go client.GetFromHistory(ctx, group, api, cfg, &waitGroup, log)
		}
		waitGroup.Wait()

		return nil
	}); err != nil {
		log.Error(err)
	}
}

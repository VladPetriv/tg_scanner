package app

import (
	"context"
	"fmt"
	"sync"

	"github.com/VladPetriv/tg_scanner/internal/client"
	"github.com/VladPetriv/tg_scanner/internal/client/auth"
	"github.com/VladPetriv/tg_scanner/internal/controller"
	"github.com/VladPetriv/tg_scanner/internal/store"
	"github.com/VladPetriv/tg_scanner/pkg/config"
	"github.com/VladPetriv/tg_scanner/pkg/logger"
	"github.com/gotd/td/telegram"
)

func Run(store *store.Store, queue controller.Controller, cfg *config.Config, log *logger.Logger) {
	tgClient, err := telegram.ClientFromEnvironment(telegram.Options{}) //nolint
	if err != nil {
		log.Error().Err(err).Msg("create telegram client")
	}

	api := tgClient.API()
	ctx := context.Background()

	appClient := client.New(ctx, store, queue, api, log, cfg)

	if err := tgClient.Run(ctx, func(ctx context.Context) error {
		var waitGroup sync.WaitGroup

		user, err := auth.Login(ctx, tgClient, cfg)
		if err != nil {
			return fmt.Errorf("authenticate user: %w", err)
		}

		userData, _ := user.GetUser().AsNotEmpty()

		groups, err := appClient.ValidateAndPushGroupsToQueue(ctx)
		if err != nil {
			return fmt.Errorf("failed to validation and push groups to queue: %w", err)
		}

		waitGroup.Add(3) //nolint

		go appClient.PushMessagesToQueue()
		go appClient.GetHistoryMessages(groups[5:])
		go appClient.GetIncomingMessages(*userData, groups[5:])

		waitGroup.Wait()

		return nil
	}); err != nil {
		log.Fatal().Err(err).Msg("start application")
	}
}

package app

import (
	"context"
	"fmt"
	"sync"

	"github.com/gotd/td/telegram"

	"github.com/VladPetriv/tg_scanner/internal/client"
	"github.com/VladPetriv/tg_scanner/internal/client/auth"
	"github.com/VladPetriv/tg_scanner/internal/controller"
	"github.com/VladPetriv/tg_scanner/internal/store"
	"github.com/VladPetriv/tg_scanner/pkg/config"
	"github.com/VladPetriv/tg_scanner/pkg/errors"
	"github.com/VladPetriv/tg_scanner/pkg/file"
	"github.com/VladPetriv/tg_scanner/pkg/logger"
)

func Run(store *store.Store, queue controller.Controller, cfg *config.Config, log *logger.Logger) {
	waitGroup := sync.WaitGroup{}

	tgClient, err := telegram.ClientFromEnvironment(telegram.Options{})
	if err != nil {
		log.Error().Err(&errors.CreateError{Name: "telegram client", ErrorValue: err}).Msg("failed to create telegram client")
	}

	api := tgClient.API()
	ctx := context.Background()

	appClient := client.New(ctx, store, queue, api, log, cfg)

	log.Info().Msg("start the application")

	if err := tgClient.Run(ctx, func(ctx context.Context) error {
		log.Info().Msg("authenticate user")

		user, err := auth.Login(ctx, tgClient, cfg)
		if err != nil {
			return fmt.Errorf("failed to authenticate user: %w", err)
		}

		userData, _ := user.GetUser().AsNotEmpty()

		log.Info().Msg("get user groups")

		groups, err := appClient.Groups.GetGroups(ctx)
		if err != nil {
			log.Error().Err(err).Msg("failed to get user groups")
		}

		log.Info().Msg("create base files")
		err = file.CreateFilesForGroups(groups)
		if err != nil {
			log.Error().Err(err).Msg("failed to create base files for group messages")
		}

		log.Info().Msg("check if groups are in cache")
		for _, groupData := range groups {
			if groupData.ID == 0 {
				continue
			}

			groupValue, err := store.Cache.Get(ctx, store.Cache.GenerateKey(groupData))
			if err != nil {
				log.Error().Err(err).Msg("failed to get value from cache")
			}

			if groupValue == "" {
				err = store.Cache.Set(ctx, store.Cache.GenerateKey(groupData), true)
				if err != nil {
					log.Error().Err(err).Msg("failed to set value into cache")
				}
			} else {
				continue
			}

			groupPhotoData, err := appClient.Groups.GetGroupPhoto(ctx, &groupData)
			if err != nil {
				log.Error().Err(err).Msgf("failed to get [%s] photo data", groupData.Username)

				continue
			}

			groupImageUrl, err := appClient.Photos.ProcessPhoto(ctx, groupPhotoData, groupData.Username)
			if err != nil {
				log.Error().Err(err).Msgf("failed to process [%s] photo data", groupData.Username)
			}

			groupData.ImageURL = groupImageUrl

			err = queue.PushDataToQueue("groups", groupData)
			if err != nil {
				log.Error().Err(err).Msgf("failed to push [%s] into queue", groupData.Username)
			}
		}

		log.Info().Msg("successfully pushed groups into queue")

		waitGroup.Add(3)

		log.Info().Msg("wait default start timeout [20s]")

		go appClient.PushToQueue()
		go appClient.GetHistoryMessages(groups[5:])
		go appClient.GetIncomingMessages(userData, groups[5:])

		waitGroup.Wait()

		return nil
	}); err != nil {
		log.Error().Err(err)
	}
}

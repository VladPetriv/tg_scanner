package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/VladPetriv/tg_scanner/config"
	"github.com/VladPetriv/tg_scanner/internal/auth"
	"github.com/VladPetriv/tg_scanner/internal/channel"
	"github.com/VladPetriv/tg_scanner/internal/file"
	"github.com/VladPetriv/tg_scanner/internal/filter"
	"github.com/VladPetriv/tg_scanner/internal/message"
	"github.com/VladPetriv/tg_scanner/logger"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
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

	// Create new client
	client, err := telegram.ClientFromEnvironment(telegram.Options{}) // nolint
	if err != nil {
		log.Errorf("ERROR_WHILE_CREATING_CLIENT:%s", err)
	}

	// Create API
	api := client.API()

	if err := client.Run(context.Background(), func(ctx context.Context) error {
		// Authorization to telegram
		user, err := auth.Login(ctx, client, cfg)
		if err != nil {
			return fmt.Errorf("AUTH_ERROR:%w", err)
		}

		waitGroup.Add(2) // nolint
		// Get user data
		u, _ := user.GetUser().AsNotEmpty()

		// Getting incoming messages
		go GetNewMessage(ctx, u, api, &waitGroup, log)

		// Getting all groups
		groups, err := channel.GetAllGroups(ctx, api)
		if err != nil {
			return fmt.Errorf("GROUPS_ERROR:%w", err)
		}

		// Create files for groups
		file.CreateFilesForGroups(groups)

		// Getting group history
		for _, group := range groups {
			go GetFromHistory(ctx, group, api, cfg, &waitGroup, log)
		}
		waitGroup.Wait()

		return nil
	}); err != nil {
		log.Error(err)
	}
}

func GetFromHistory(ctx context.Context, group channel.Group, api *tg.Client, cfg *config.Config, wg *sync.WaitGroup, log *logger.Logger) { // nolint
	defer wg.Done()

	path := fmt.Sprintf("./data/%s.json", group.Username)

	log.Info("Start parsing messages from telgram")
	for {
		data, err := channel.GetChannelHistory(ctx, cfg.Limit, tg.InputPeerChannel{
			ChannelID:  int64(group.ID),
			AccessHash: int64(group.AccessHash),
		}, api)
		if err != nil {
			log.Error(err)
		}

		modifiedData, _ := data.AsModified()

		messages := message.GetMessagesFromTelegram(ctx, modifiedData, &tg.InputPeerChannel{
			ChannelID:  int64(group.ID),
			AccessHash: int64(group.AccessHash),
		}, api)

		messagesFromFile, err := file.GetMessagesFromFile(path)
		if err != nil {
			log.Error(err)
		}

		for index := range messages {
			msg, ok := filter.Messages(&messages[index])
			if !ok {
				continue
			}

			messagesFromFile = append(messagesFromFile, *msg)
		}

		result := filter.RemoveDuplicateByMessage(messagesFromFile)

		err = file.WriteMessagesToFile(result, path)
		if err != nil {
			log.Error(err)
		}

		time.Sleep(time.Minute * 2)
	}
}

func GetNewMessage(ctx context.Context, user *tg.User, api *tg.Client, wg *sync.WaitGroup, log *logger.Logger) {
	defer wg.Done()

	path := "./data/incoming.json"

	err := file.CreateFileForIncoming()
	if err != nil {
		log.Error(err)
	}

	log.Info("Start getting incoming messages")
	for {
		messagesFromFile, err := file.GetMessagesFromFile(path)
		if err != nil {
			log.Error(err)
		}

		incomingMessage, err := message.GetIncomingMessages(ctx, user, api)
		if err != nil {
			log.Error(err)
		}

		for index := range incomingMessage {
			msg, ok := filter.Messages(&incomingMessage[index])
			if !ok {
				continue
			}

			messagesFromFile = append(messagesFromFile, *msg)
		}

		result := filter.RemoveDuplicateByMessage(messagesFromFile)

		err = file.WriteMessagesToFile(result, path)
		if err != nil {
			log.Error(err)
		}

		time.Sleep(time.Minute) // nolint
	}
}

package client

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
	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/internal/service"
	"github.com/VladPetriv/tg_scanner/logger"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
)

func GetFromHistory(ctx context.Context, group channel.Group, api *tg.Client, cfg *config.Config, wg *sync.WaitGroup, log *logger.Logger) { // nolint
	time.Sleep(time.Second * 10)
	defer wg.Done()

	path := fmt.Sprintf("./data/%s.json", group.Username)

	for {
		log.Info("Start parsing messages from telgram")
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
			msg.PeerID = group
			messagesFromFile = append(messagesFromFile, *msg)
		}

		result := filter.RemoveDuplicateByMessage(messagesFromFile)

		err = file.WriteMessagesToFile(result, path)
		if err != nil {
			log.Error(err)
		}

		time.Sleep(time.Minute * 10)
	}
}

func GetNewMessage(ctx context.Context, user *tg.User, api *tg.Client, groups []channel.Group, wg *sync.WaitGroup, log *logger.Logger) {
	defer wg.Done()

	path := "./data/incoming.json"

	err := file.CreateFileForIncoming()
	if err != nil {
		log.Error(err)
	}

	for {
		log.Info("Start getting incoming messages")
		messagesFromFile, err := file.GetMessagesFromFile(path)
		if err != nil {
			log.Error(err)
		}

		incomingMessage, err := message.GetIncomingMessages(ctx, user, groups, api)
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

func SaveToDb(ctx context.Context, serviceManager *service.Manager, api *tg.Client, log *logger.Logger) {
	for {
		log.Info("Start saving messages to db")
		messages, err := file.ParseFromFiles("data")
		if err != nil {
			log.Error(err)
		}

		for _, msg := range messages {
			err := message.GetRepliesForMessageBeforeSave(ctx, &msg, api)
			if err != nil {
				log.Error(err)
			}

			channel, err := serviceManager.Channel.GetChannelByName(msg.PeerID.Username)
			if err != nil {
				log.Error(err)
			}

			err = serviceManager.Message.CreateMessage(&model.Message{ChannelID: channel.ID, Title: msg.Message})
			if err != nil {
				log.Error(err)
			}

			for _, replie := range msg.Replies.Messages {
				replieCandidate, err := serviceManager.Replie.GetReplieByName(replie.Message)
				if replieCandidate != nil || err != nil {
					log.Error(err)
					continue
				}

				message, err := serviceManager.Message.GetMessageByName(msg.Message)
				if err != nil {
					log.Error(err)
				}

				err = serviceManager.Replie.CreateReplie(&model.Replie{MessageID: message.ID, Title: replie.Message})
				if err != nil {
					log.Error(err)
				}
			}
		}
		time.Sleep(time.Minute * 15)
	}
}

func Run(serviceManager *service.Manager, waitGroup *sync.WaitGroup, cfg *config.Config, log *logger.Logger) {
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

		waitGroup.Add(3) // nolint

		// Get user data
		uData, _ := user.GetUser().AsNotEmpty()

		// Getting all groups
		groups, err := channel.GetAllGroups(ctx, api)
		if err != nil {
			return fmt.Errorf("GROUPS_ERROR:%w", err)
		}

		// Getting incoming messages
		go GetNewMessage(ctx, uData, api, groups, waitGroup, log)

		// Create files for groups
		file.CreateFilesForGroups(groups)

		// Getting group history
		for _, group := range groups {
			err := serviceManager.Channel.CreateChannel(&model.Channel{Name: group.Username})
			if err != nil {
				log.Error(err)
			}

			// go GetFromHistory(ctx, group, api, cfg, waitGroup, log)
		}

		time.Sleep(time.Second * 5)
		go SaveToDb(ctx, serviceManager, api, log)

		waitGroup.Wait()
		return nil
	}); err != nil {
		log.Error(err)
	}
}

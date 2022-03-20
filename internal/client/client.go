package client

import (
	"context"
	"fmt"
	"github.com/VladPetriv/tg_scanner/config"
	"github.com/VladPetriv/tg_scanner/internal/channel"
	"github.com/VladPetriv/tg_scanner/internal/file"
	"github.com/VladPetriv/tg_scanner/internal/filter"
	"github.com/VladPetriv/tg_scanner/internal/message"
	"github.com/VladPetriv/tg_scanner/logger"
	"github.com/gotd/td/tg"
	"sync"
	"time"
)

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

		time.Sleep(time.Minute * 5)
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

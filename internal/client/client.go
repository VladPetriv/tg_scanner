package client

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/VladPetriv/tg_scanner/config"
	"github.com/VladPetriv/tg_scanner/internal/auth"
	"github.com/VladPetriv/tg_scanner/internal/channel"
	"github.com/VladPetriv/tg_scanner/internal/file"
	"github.com/VladPetriv/tg_scanner/internal/filter"
	"github.com/VladPetriv/tg_scanner/internal/firebase"
	"github.com/VladPetriv/tg_scanner/internal/message"
	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/internal/service"
	"github.com/VladPetriv/tg_scanner/internal/user"
	"github.com/VladPetriv/tg_scanner/logger"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
)

func GetMessagesFromHistory(ctx context.Context, channels []channel.Channel, wg *sync.WaitGroup, api *tg.Client, log *logger.Logger) {
	time.Sleep(time.Second * 20)
	defer wg.Done()
	for {
		for _, chnl := range channels {
			log.Infof("Start getting messages from history[%s]", chnl.Username)
			fileName := fmt.Sprintf("./data/%s.json", chnl.Username)

			data, err := channel.GetChannelHistory(ctx, &tg.InputPeerChannel{
				ChannelID:  int64(chnl.ID),
				AccessHash: int64(chnl.AccessHash),
			}, api)
			if err != nil {
				log.Error(err)
			}

			modifiedData, _ := data.AsModified()

			messages := message.GetMessagesFromTelegram(ctx, modifiedData, &tg.InputPeerChannel{
				ChannelID:  int64(chnl.ID),
				AccessHash: int64(chnl.AccessHash),
			}, api)

			messagesFromFile, err := file.GetMessagesFromFile(fileName)
			if err != nil {
				log.Error(err)
			}

			for _, msg := range messages {
				msg, ok := filter.Messages(&msg)
				if !ok {
					continue
				}

				msg.PeerID = chnl

				u, err := user.GetUserInfo(ctx, msg.FromID.UserID, msg.ID, &tg.InputPeerChannel{
					ChannelID:  int64(chnl.ID),
					AccessHash: int64(chnl.AccessHash),
				}, api)
				if err != nil {
					log.Error(err)

					continue
				}

				msg.FromID = *u
				messagesFromFile = append(messagesFromFile, *msg)
			}

			result := filter.RemoveDuplicateByMessage(messagesFromFile)

			err = file.WriteMessagesToFile(result, fileName)
			if err != nil {
				log.Error(err)
			}

			time.Sleep(time.Second * 10)
		}

		time.Sleep(time.Minute * 30)
	}
}

func GetNewMessage(ctx context.Context, user *tg.User, api *tg.Client, channels []channel.Channel, wg *sync.WaitGroup, log *logger.Logger) {
	defer wg.Done()
	time.Sleep(time.Second * 20)

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

		incomingMessage, err := message.GetIncomingMessages(ctx, user, channels, api)
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

func SaveToDb(ctx context.Context, serviceManager *service.Manager, cfg *config.Config, api *tg.Client, log *logger.Logger) {
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

			fileName, err := user.ProcessUserPhoto(ctx, &msg.FromID, api)
			if err != nil {
				log.Error(err)
			}

			userImageUrl, err := firebase.SendImageToStorage(ctx, cfg, fileName, msg.FromID.Username)
			if err != nil {
				log.Error(err)
			}

			fullName := fmt.Sprintf("%s %s", msg.FromID.FirstName, msg.FromID.LastName)
			userID, err := serviceManager.User.CreateUser(&model.User{Username: msg.FromID.Username, FullName: fullName, PhotoURL: userImageUrl})
			if err != nil {
				log.Error(err)
			}

			messageID, err := serviceManager.Message.CreateMessage(&model.Message{ChannelID: channel.ID, UserID: userID, Title: msg.Message})
			if err != nil {
				log.Error(err)
			}

			for _, replie := range msg.Replies.Messages {
				fileName, err := user.ProcessUserPhoto(ctx, &replie.FromID, api)
				if err != nil {
					log.Error(err)
				}

				userImageUrl, err := firebase.SendImageToStorage(ctx, cfg, fileName, replie.FromID.Username)
				if err != nil {
					log.Error(err)
				}

				fullName := fmt.Sprintf("%s %s", replie.FromID.FirstName, replie.FromID.LastName)
				userID, err := serviceManager.User.CreateUser(&model.User{Username: replie.FromID.Username, FullName: fullName, PhotoURL: userImageUrl})
				if err != nil {
					log.Error(err)
				}

				err = serviceManager.Replie.CreateReplie(&model.Replie{UserID: userID, MessageID: messageID, Title: replie.Message})
				if err != nil {
					log.Error(err)
				}
			}
		}

		time.Sleep(time.Minute * 15)
	}
}

func RemoveMessageWithOutReplies(serviceManager *service.Manager, log *logger.Logger) {
	for {
		log.Infof("Start remove messages without replies")
		messages, err := serviceManager.Message.GetMessagesWithRepliesCount()
		if err != nil {
			log.Error(err)
		}

		for _, message := range messages {
			if message.RepliesCount == 0 {
				err := serviceManager.Message.DeleteMessageByID(message.ID)
				if err != nil {
					log.Error(err)
				}
				continue
			}

			continue
		}

		time.Sleep(time.Minute * 60)
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

		// Getting all channel
		channels, err := channel.GetAllChannels(ctx, api)
		if err != nil {
			return fmt.Errorf("Channel_ERROR:%w", err)
		}

		err = file.CreateFilesForChannels(channels)
		if err != nil {
			log.Error(err)
		}

		// Getting channel history
		for _, chnl := range channels {
			candidate, err := serviceManager.Channel.GetChannelByName(chnl.Username)
			if err != nil {
				log.Error(err)
			}
			if candidate != nil {
				continue
			}

			filename, err := channel.ProcessChannelPhoto(ctx, &chnl, api)
			if err != nil {
				log.Error(err)
			}

			channelImageURL, err := firebase.SendImageToStorage(ctx, cfg, filename, chnl.Username)
			if err != nil {
				log.Error(err)
			}

			err = serviceManager.Channel.CreateChannel(&model.Channel{Name: chnl.Username, Title: chnl.Title, PhotoURL: channelImageURL})
			if err != nil {
				log.Error(err)
			}

			os.Remove(fmt.Sprintf("images/%s.jpg", chnl.Username))
		}

		go SaveToDb(ctx, serviceManager, cfg, api, log)

		go GetNewMessage(ctx, uData, api, channels, waitGroup, log)
		go GetMessagesFromHistory(ctx, channels, waitGroup, api, log)
		go RemoveMessageWithOutReplies(serviceManager, log)

		waitGroup.Wait()

		return nil
	}); err != nil {
		log.Error(err)
	}
}

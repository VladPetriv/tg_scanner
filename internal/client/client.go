package client

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"

	"github.com/VladPetriv/tg_scanner/internal/client/auth"
	"github.com/VladPetriv/tg_scanner/internal/client/channel"
	"github.com/VladPetriv/tg_scanner/internal/client/message"
	"github.com/VladPetriv/tg_scanner/internal/client/photo"
	"github.com/VladPetriv/tg_scanner/internal/client/replie"
	"github.com/VladPetriv/tg_scanner/internal/client/user"
	"github.com/VladPetriv/tg_scanner/internal/file"
	"github.com/VladPetriv/tg_scanner/internal/filter"
	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/internal/service"
	"github.com/VladPetriv/tg_scanner/internal/store/redis"
	"github.com/VladPetriv/tg_scanner/pkg/config"
	"github.com/VladPetriv/tg_scanner/pkg/logger"
	"github.com/VladPetriv/tg_scanner/pkg/utils"
)

// Timeouts
var (
	startTimeout    time.Duration = 20 * time.Second
	historyTimeout  time.Duration = 30 * time.Minute
	removeTimeout   time.Duration = 30 * time.Minute
	saveTimeout     time.Duration = 15 * time.Minute
	incomingTimeout time.Duration = time.Minute
)

func GetMessagesFromHistory(ctx context.Context, channels []model.TgChannel, api *tg.Client, log *logger.Logger) {
	time.Sleep(startTimeout)

	for {
		for _, channelData := range channels {
			log.Infof("Start getting messages from history[%s]", channelData.Username)
			fileName := fmt.Sprintf("./data/%s.json", channelData.Username)

			data, err := channel.GetChannelHistory(ctx, &tg.InputPeerChannel{
				ChannelID:  channelData.ID,
				AccessHash: channelData.AccessHash,
			}, api)
			if err != nil {
				log.Error(err)
			}

			modifiedData, _ := data.AsModified()

			messages := message.GetMessagesFromTelegram(ctx, modifiedData, &tg.InputPeerChannel{
				ChannelID:  channelData.ID,
				AccessHash: channelData.AccessHash,
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

				msg.PeerID = channelData

				userInfo, err := user.GetUserInfo(ctx, msg.FromID.UserID, msg.ID, &tg.InputPeerChannel{
					ChannelID:  channelData.ID,
					AccessHash: channelData.AccessHash,
				}, api)
				if err != nil {
					log.Error(err)

					continue
				}

				msg.FromID = *userInfo
				messagesFromFile = append(messagesFromFile, *msg)
			}

			result := filter.RemoveDuplicateInMessage(messagesFromFile)

			err = file.WriteMessagesToFile(result, fileName)
			if err != nil {
				log.Error(err)
			}

			time.Sleep(time.Second * 10)
		}

		time.Sleep(historyTimeout)
	}
}

func GetNewMessage(ctx context.Context, user *tg.User, api *tg.Client, channels []model.TgChannel, log *logger.Logger) {
	time.Sleep(startTimeout)

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

		result := filter.RemoveDuplicateInMessage(messagesFromFile)

		err = file.WriteMessagesToFile(result, path)
		if err != nil {
			log.Error(err)
		}

		time.Sleep(incomingTimeout)
	}
}

func SaveToDb(ctx context.Context, serviceManager *service.Manager, redisDB *redis.RedisDB, cfg *config.Config, api *tg.Client, log *logger.Logger) {
	for {
		log.Info("Start saving messages to db")

		messages, err := file.ParseFromFiles("data")
		if err != nil {
			log.Error(err)
		}

		for _, messageData := range messages {
			messageValue, err := redisDB.GetMessageFromRedis(ctx, redis.GenerateMessageKey(messageData))
			if err != nil {
				log.Error(err)
			}

			if messageValue == "" {
				err := redisDB.SetMessageToRedis(ctx, redis.GenerateMessageKey(messageData), true)
				if err != nil {
					log.Error(err)
				}
			} else {
				continue
			}

			err = replie.GetRepliesForMessageBeforeSave(ctx, &messageData, api)
			if err != nil {
				log.Error(err)
			}

			filter.RemoveDuplicateInReplies(&messageData.Replies)

			channel, _ := serviceManager.Channel.GetChannelByName(messageData.PeerID.Username)

			userPhotoData, err := user.GetUserPhoto(ctx, messageData.FromID, api)
			if err != nil {
				log.Error(err)
			}

			userImageUrl, err := photo.ProcessPhoto(ctx, userPhotoData, messageData.FromID.Username, cfg, api)
			if err != nil {
				log.Error(err)
			}

			userID, err := serviceManager.User.CreateUser(&model.User{
				Username: messageData.FromID.Username,
				FullName: fmt.Sprintf("%s %s", messageData.FromID.FirstName, messageData.FromID.LastName),
				ImageURL: userImageUrl,
			})
			if _, ok := err.(*utils.RecordIsExistError); ok && err != nil {
				log.Warn(err)
			} else if err != nil {
				log.Error(err)
			}

			var messageImageUrl string

			if ok, _ := message.CheckMessagePhotoStatus(ctx, &messageData, api); ok {
				messagePhotoData, err := message.GetMessagePhoto(ctx, messageData, api)
				if err != nil {
					log.Error(err)
				}

				messageImageUrl, err = photo.ProcessPhoto(ctx, messagePhotoData, fmt.Sprint(messageData.ID), cfg, api)
				if err != nil {
					log.Error(err)
				}
			}

			messageID, err := serviceManager.Message.CreateMessage(&model.Message{
				ChannelID:  channel.ID,
				UserID:     userID,
				Title:      messageData.Message,
				MessageURL: fmt.Sprintf("https://t.me/%s/%d", messageData.PeerID.Username, messageData.ID),
				ImageURL:   messageImageUrl,
			})
			if _, ok := err.(*utils.RecordIsExistError); ok && err != nil {
				log.Warn(err)
			} else if err != nil {
				log.Error(err)
			}

			for _, replieData := range messageData.Replies.Messages {
				userPhotoData, err := user.GetUserPhoto(ctx, replieData.FromID, api)
				if err != nil {
					log.Error(err)
				}

				userImageUrl, err := photo.ProcessPhoto(ctx, userPhotoData, replieData.FromID.Username, cfg, api)
				if err != nil {
					log.Error(err)
				}

				userID, err := serviceManager.User.CreateUser(&model.User{
					Username: replieData.FromID.Username,
					FullName: fmt.Sprintf("%s %s", replieData.FromID.FirstName, replieData.FromID.LastName),
					ImageURL: userImageUrl,
				})
				if _, ok := err.(*utils.RecordIsExistError); ok && err != nil {
					log.Warn(err)
				} else if err != nil {
					log.Error(err)
				}

				var replieImageUrl string

				if replieData.Media.Photo != nil {
					repliePhotoData, err := replie.GetRepliePhoto(ctx, replieData, api)
					if err != nil {
						log.Error(err)
					}

					replieImageUrl, err = photo.ProcessPhoto(ctx, repliePhotoData, replieData.ID, cfg, api)
					if err != nil {
						log.Error(err)
					}
				}

				err = serviceManager.Replie.CreateReplie(&model.Replie{
					UserID:    userID,
					MessageID: messageID,
					Title:     replieData.Message,
					ImageURL:  replieImageUrl,
				})
				if _, ok := err.(*utils.RecordIsExistError); ok && err != nil {
					log.Warn(err)
				} else if err != nil {
					log.Error(err)
				}
			}
		}

		time.Sleep(saveTimeout)
	}
}

func RemoveMessageWithOutReplies(serviceManager *service.Manager, log *logger.Logger) {
	for {
		log.Infof("Start remove messages without replies")

		messages, err := serviceManager.Message.GetMessagesWithRepliesCount()
		if err != nil {
			log.Warn(err)
		}

		for _, message := range messages {
			if message.RepliesCount == 0 {
				err := serviceManager.Message.DeleteMessageByID(message.ID)
				if err != nil {
					log.Warn(err)
				}

				continue
			}

			continue
		}

		time.Sleep(removeTimeout)
	}
}

func Run(serviceManager *service.Manager, redisDB *redis.RedisDB, waitGroup *sync.WaitGroup, cfg *config.Config, log *logger.Logger) {
	tgClient, err := telegram.ClientFromEnvironment(telegram.Options{}) // nolint
	if err != nil {
		log.Error(&utils.CreateError{Name: "telegram client", ErrorValue: err})
	}

	api := tgClient.API()

	if err := tgClient.Run(context.Background(), func(ctx context.Context) error {
		user, err := auth.Login(ctx, tgClient, cfg)
		if err != nil {
			return fmt.Errorf("AUTH_ERROR:%w", err)
		}

		waitGroup.Add(4)

		uData, _ := user.GetUser().AsNotEmpty()

		channels, err := channel.GetAllChannels(ctx, api)
		if err != nil {
			log.Error(err)
		}

		err = file.CreateFilesForChannels(channels)
		if err != nil {
			log.Error(err)
		}

		for _, channelData := range channels {
			candidate, _ := serviceManager.Channel.GetChannelByName(channelData.Username)
			if candidate != nil {
				continue
			}

			channelPhotoData, err := channel.GetChannelPhoto(ctx, &channelData, api)
			if err != nil {
				log.Error(err)
			}

			channelImageUrl, err := photo.ProcessPhoto(ctx, channelPhotoData, channelData.Username, cfg, api)
			if err != nil {
				log.Error(err)
			}

			err = serviceManager.Channel.CreateChannel(&model.Channel{
				Name:     channelData.Username,
				Title:    channelData.Title,
				ImageURL: channelImageUrl,
			})
			if err != nil {
				log.Error(err)
			}
		}

		go SaveToDb(ctx, serviceManager, redisDB, cfg, api, log)
		go GetNewMessage(ctx, uData, api, channels, log)
		go GetMessagesFromHistory(ctx, channels, api, log)
		go RemoveMessageWithOutReplies(serviceManager, log)

		waitGroup.Wait()

		return nil
	}); err != nil {
		log.Error(err)
	}
}

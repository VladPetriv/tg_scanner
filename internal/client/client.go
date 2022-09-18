package client

import (
	"context"
	"fmt"
	"time"

	"github.com/gotd/td/tg"
	"github.com/rs/zerolog/log"

	"github.com/VladPetriv/tg_scanner/internal/client/group"
	"github.com/VladPetriv/tg_scanner/internal/client/message"
	"github.com/VladPetriv/tg_scanner/internal/client/photo"
	"github.com/VladPetriv/tg_scanner/internal/client/reply"
	"github.com/VladPetriv/tg_scanner/internal/client/user"
	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/internal/store"
	"github.com/VladPetriv/tg_scanner/pkg/file"
	"github.com/VladPetriv/tg_scanner/pkg/filter"
	"github.com/VladPetriv/tg_scanner/pkg/logger"
)

// Timeouts
var (
	_startTimeout    time.Duration = 20 * time.Second
	_historyTimeout  time.Duration = 30 * time.Minute
	_removeTimeout   time.Duration = 30 * time.Minute
	_saveTimeout     time.Duration = 15 * time.Minute
	_incomingTimeout time.Duration = time.Minute
)

type appClient struct {
	ctx   context.Context
	store *store.Store
	api   *tg.Client
	log   *logger.Logger

	Groups   group.Group
	Messages message.Message
	Users    user.User
	Photo    photo.Photo
	Replies  reply.Reply
}

var _ AppClient = (*appClient)(nil)

func New(ctx context.Context, store *store.Store, api *tg.Client, log *logger.Logger) *appClient {
	return &appClient{
		ctx:      ctx,
		store:    store,
		api:      api,
		log:      log,
		Groups:   group.New(log, api),
		Messages: message.New(log, api),
		Users:    user.New(log, api),
		Replies:  reply.New(log, api),
	}
}

func (c appClient) GetHistoryMessages(groups []model.TgGroup) {
	time.Sleep(_startTimeout)

	for {
		for _, groupData := range groups {
			c.log.Info().Msgf("get - [%s]", groupData.Username)

			path := fmt.Sprintf("./data/%s.json", groupData.Username)

			groupMessages, err := c.Groups.GetMessagesFromGroupHistory(c.ctx, &tg.InputPeerChannel{
				ChannelID:  groupData.ID,
				AccessHash: groupData.AccessHash,
			})
			if err != nil {
				c.log.Error().Err(err).Msg("failed to get messages from group history")
			}

			modifiedGroupMessages, ok := groupMessages.AsModified()
			if !ok {
				c.log.Warn().Msg("failed to get modified group messages")
			}

			processedMessages := c.Messages.ProcessHistoryMessages(c.ctx, modifiedGroupMessages, &tg.InputPeerChannel{
				ChannelID:  groupData.ID,
				AccessHash: groupData.AccessHash,
			})

			messagesFromFile, err := file.GetMessagesFromFile(path)
			if err != nil {
				c.log.Error().Err(err).Msg("failed to get messages from file")
			}

			for _, msg := range processedMessages {
				// check if message is question
				ok := filter.Message(&msg)
				if !ok {
					continue
				}

				msg.PeerID = groupData

				// get user info for message
				userInfo, err := c.Users.GetUser(c.ctx, msg, &tg.InputPeerChannel{
					ChannelID:  groupData.ID,
					AccessHash: groupData.AccessHash,
				})
				if err != nil {
					c.log.Error().Err(err).Msg("failed to get user info for message")

					continue
				}

				msg.FromID = *userInfo

				// get replies for message
				replies, err := c.Replies.GetReplies(c.ctx, &msg, &tg.InputPeerChannel{
					ChannelID:  groupData.ID,
					AccessHash: groupData.AccessHash,
				})
				if err != nil {
					c.log.Error().Err(err).Msg("failed to get replies for message")

					continue
				}

				processedReplies := c.Replies.ProcessReplies(c.ctx, replies, &tg.InputPeerChannel{
					ChannelID:  groupData.ID,
					AccessHash: groupData.AccessHash,
				})

				// get user info for replies
				for _, reply := range processedReplies {
					userInfo, err := c.Users.GetUser(c.ctx, reply, &tg.InputPeerChannel{
						ChannelID:  groupData.ID,
						AccessHash: groupData.AccessHash,
					})
					if err != nil {
						c.log.Error().Err(err).Msg("failed to get user info for reply")

						continue
					}

					reply.FromID = *userInfo
				}

				msg.Replies.Count = len(processedReplies)
				msg.Replies.Messages = processedReplies
			}

			messagesFromFile = append(messagesFromFile, processedMessages...)

			result := filter.RemoveDuplicatesFromMessages(messagesFromFile)

			err = file.WriteMessagesToFile(result, path)
			if err != nil {
				c.log.Error().Err(err).Msg("failed to write messages into file")
			}

			time.Sleep(time.Second * 10)
		}

		time.Sleep(_historyTimeout)
	}
}

func (c appClient) GetIncomingMessages(user *tg.User, groups []model.TgGroup) {
	time.Sleep(_startTimeout)

	path := "./data/incoming.json"

	c.log.Info().Msg("create base file for incoming messages")
	err := file.CreateFileForIncoming()
	if err != nil {
		c.log.Error().Err(err).Msg("failed to create base file for incoming messages")
	}

	for {
		c.log.Info().Msg("get - [incoming messages]")

		processedMessages, err := c.Messages.ProcessIncomingMessages(c.ctx, user, groups)
		if err != nil {
			log.Error().Err(err).Msg("failed to process incoming messages")
		}

		messagesFromFile, err := file.GetMessagesFromFile(path)
		if err != nil {
			log.Error().Err(err).Msg("failed to get message from files")
		}

		for _, msg := range processedMessages {
			// check if message in question
			ok := filter.Message(&msg)
			if !ok {
				continue
			}

			// get user info for message
			userInfo, err := c.Users.GetUser(c.ctx, msg, &tg.InputPeerChannel{
				ChannelID:  msg.PeerID.ID,
				AccessHash: msg.PeerID.AccessHash,
			})
			if err != nil {
				c.log.Error().Err(err).Msg("failed to get user info for message")

				continue
			}

			msg.FromID = *userInfo
		}

		messagesFromFile = append(messagesFromFile, processedMessages...)

		result := filter.RemoveDuplicatesFromMessages(messagesFromFile)

		err = file.WriteMessagesToFile(result, path)
		if err != nil {
			log.Error().Err(err).Msg("failed to write messages intofile")
		}

		time.Sleep(_incomingTimeout)
	}
}

//TODO: refactor it
/*
func SaveToKafka(ctx context.Context, redisDB *redis.RedisDB, cfg *config.Config, api *tg.Client, log *logger.Logger) {
	for {
		log.Info("Start saving messages to db")

		messages, err := file.ParseFromFiles("data")
		if err != nil {
			log.Error(err)
		}

		for _, messageData := range messages {
			messageValue, err := redisDB.GetDataFromRedis(ctx, redis.GenerateKey(messageData))
			if err != nil {
				log.Error(err)
			}

			if messageValue == "" {
				err := redisDB.SetDataToRedis(ctx, redis.GenerateKey(messageData), true)
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

			if len(messageData.Replies.Messages) == 0 {
				continue
			}

			userPhotoData, err := user.GetUserPhoto(ctx, messageData.FromID, api)
			if err != nil {
				log.Error(err)
			}

			userImageUrl, err := photo.ProcessPhoto(ctx, userPhotoData, messageData.FromID.Username, cfg, api)
			if err != nil {
				log.Error(err)
			}

			messageData.FromID.ImageURL = userImageUrl
			messageData.FromID.Fullname = fmt.Sprintf("%s %s", messageData.FromID.FirstName, messageData.FromID.LastName)

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

			messageData.MessageURL = fmt.Sprintf("https://t.me/%s/%d", messageData.PeerID.Username, messageData.ID)
			messageData.ImageURL = messageImageUrl

			for index, replieData := range messageData.Replies.Messages {
				userPhotoData, err := user.GetUserPhoto(ctx, replieData.FromID, api)
				if err != nil {
					log.Error(err)
				}

				userImageUrl, err := photo.ProcessPhoto(ctx, userPhotoData, replieData.FromID.Username, cfg, api)
				if err != nil {
					log.Error(err)
				}

				messageData.Replies.Messages[index].FromID.ImageURL = userImageUrl
				messageData.Replies.Messages[index].FromID.Fullname = fmt.Sprintf("%s %s", replieData.FromID.FirstName, replieData.FromID.LastName)

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

				replieData.ImageURL = replieImageUrl
			}

			err = kafka.PushDataToQueue("messages.get", cfg.KafkaAddr, messageData)
			if err != nil {
				log.Error(err)
			}
		}

		time.Sleep(saveTimeout)
	}
}
*/

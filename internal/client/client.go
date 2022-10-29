package client

import (
	"context"
	"fmt"
	"time"

	"github.com/VladPetriv/tg_scanner/internal/client/group"
	"github.com/VladPetriv/tg_scanner/internal/client/message"
	"github.com/VladPetriv/tg_scanner/internal/client/photo"
	"github.com/VladPetriv/tg_scanner/internal/client/reply"
	"github.com/VladPetriv/tg_scanner/internal/client/user"
	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/internal/store"
	"github.com/VladPetriv/tg_scanner/pkg/config"
	"github.com/VladPetriv/tg_scanner/pkg/file"
	"github.com/VladPetriv/tg_scanner/pkg/filter"
	"github.com/VladPetriv/tg_scanner/pkg/kafka"
	"github.com/VladPetriv/tg_scanner/pkg/logger"
	"github.com/gotd/td/tg"
	"github.com/rs/zerolog/log"
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
	cfg   *config.Config

	Groups   group.Group
	Messages message.Message
	Users    user.User
	Photos   photo.Photo
	Replies  reply.Reply
}

var _ AppClient = (*appClient)(nil)

func New(ctx context.Context, store *store.Store, api *tg.Client, log *logger.Logger, cfg *config.Config) *appClient {
	return &appClient{
		ctx:      ctx,
		store:    store,
		api:      api,
		log:      log,
		cfg:      cfg,
		Groups:   group.New(log, api),
		Messages: message.New(log, api),
		Users:    user.New(log, api),
		Replies:  reply.New(log, api),
		Photos:   photo.New(log, store),
	}
}

func (c appClient) GetHistoryMessages(groups []model.TgGroup) {
	time.Sleep(_startTimeout)

	for {
		for _, groupData := range groups {
			c.log.Info().Msgf("get - [%s]", groupData.Username)

			path := fmt.Sprintf("./data/%s.json", groupData.Username)

			groupPeer := &tg.InputPeerChannel{
				ChannelID:  groupData.ID,
				AccessHash: groupData.AccessHash,
			}

			groupMessages, err := c.Groups.GetMessagesFromGroupHistory(c.ctx, groupPeer)
			if err != nil {
				c.log.Error().Err(err).Msg("failed to get messages from group history")
			}

			modifiedGroupMessages, ok := groupMessages.AsModified()
			if !ok {
				c.log.Warn().Msg("failed to get modified group messages")
			}

			processedMessages := c.Messages.ProcessHistoryMessages(c.ctx, modifiedGroupMessages, groupPeer)

			messagesFromFile, err := file.GetMessagesFromFile(path)
			if err != nil {
				c.log.Error().Err(err).Msg("failed to get messages from file")
			}

			for _, msg := range processedMessages {
				// check if message is question
				ok := filter.ProcessMessage(&msg)
				if !ok {
					continue
				}

				msg.PeerID = groupData

				// get user info for message
				userInfo, err := c.Users.GetUser(c.ctx, msg, groupPeer)
				if err != nil {
					c.log.Error().Err(err).Msg("failed to get user info for message")

					continue
				}

				msg.FromID = *userInfo

				// get replies for message
				replies, err := c.Replies.GetReplies(c.ctx, &msg, groupPeer)
				if err != nil {
					c.log.Error().Err(err).Msg("failed to get replies for message")

					continue
				}

				processedReplies := c.Replies.ProcessReplies(c.ctx, replies, groupPeer)

				// get user info for replies
				for index, reply := range processedReplies {
					userInfo, err := c.Users.GetUser(c.ctx, reply, groupPeer)
					if err != nil {
						c.log.Error().Err(err).Msg("failed to get user info for reply")

						continue
					}

					processedReplies[index].FromID = *userInfo
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
			ok := filter.ProcessMessage(&msg)
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
			log.Error().Err(err).Msg("failed to write messages into file")
		}

		time.Sleep(_incomingTimeout)
	}
}

func (c appClient) PushToQueue() {
	for {
		c.log.Info().Msg("pushing messages to queue")

		messages, err := file.ParseFromFiles("data")
		if err != nil {
			c.log.Error().Err(err).Msg("failed to get messages from files")
		}

		for _, messageData := range messages {
			messageValue, err := c.store.Cache.Get(c.ctx, c.store.Cache.GenerateKey(messageData))
			if err != nil {
				c.log.Error().Err(err).Msg("failed to get message key from cache")
			}

			if messageValue == "" {
				err := c.store.Cache.Set(c.ctx, c.store.Cache.GenerateKey(messageData), true)
				if err != nil {
					c.log.Error().Err(err).Msg("failed to set message into cache")
				}
			} else {
				continue
			}

			groupPeer := &tg.InputPeerChannel{
				ChannelID:  messageData.PeerID.ID,
				AccessHash: messageData.PeerID.AccessHash,
			}

			// get replies for message
			replies, err := c.Replies.GetReplies(c.ctx, &messageData, groupPeer)
			if err != nil {
				c.log.Error().Err(err).Msg("failed to get replies for message")

				continue
			}

			processedReplies := c.Replies.ProcessReplies(c.ctx, replies, groupPeer)

			// get user info for replies
			for _, reply := range processedReplies {
				userInfo, err := c.Users.GetUser(c.ctx, reply, groupPeer)
				if err != nil {
					c.log.Error().Err(err).Msg("failed to get user info for reply")

					continue
				}

				reply.FromID = *userInfo
			}

			messageData.Replies.Count = len(processedReplies)
			messageData.Replies.Messages = processedReplies

			filter.RemoveDuplicatesFromReplies(&messageData.Replies)

			// if len of replies is 0 move to other message
			if len(messageData.Replies.Messages) == 0 {
				continue
			}

			// process user photo
			userPhotoData, err := c.Users.GetUserPhoto(c.ctx, messageData.FromID)
			if err != nil {
				c.log.Error().Err(err).Msg("failed to get user photo")
			}

			userImageUrl, err := c.Photos.ProcessPhoto(c.ctx, userPhotoData, messageData.FromID.Username)
			if err != nil {
				c.log.Error().Err(err).Msg("failed to process user photo")
			}

			messageData.FromID.ImageURL = userImageUrl
			messageData.FromID.Fullname = fmt.Sprintf("%s %s", messageData.FromID.FirstName, messageData.FromID.LastName)

			var messageImageUrl string

			if ok, _ := c.Messages.CheckMessagePhotoStatus(c.ctx, &messageData); ok {
				messagePhotoData, err := c.Messages.GetMessagePhoto(c.ctx, messageData)
				if err != nil {
					c.log.Error().Err(err).Msg("failed to check message photo status")
				}

				messageImageUrl, err = c.Photos.ProcessPhoto(c.ctx, messagePhotoData, fmt.Sprint(messageData.ID))
				if err != nil {
					c.log.Error().Err(err).Msg("failed to process message photo")
				}
			}

			messageData.MessageURL = fmt.Sprintf("https://t.me/%s/%d", messageData.PeerID.Username, messageData.ID)
			messageData.ImageURL = messageImageUrl

			for index, replyData := range messageData.Replies.Messages {
				userPhotoData, err := c.Users.GetUserPhoto(c.ctx, replyData.FromID)
				if err != nil {
					c.log.Error().Err(err).Msg("failed to get user photo")
				}

				userImageUrl, err := c.Photos.ProcessPhoto(c.ctx, userPhotoData, replyData.FromID.Username)
				if err != nil {
					log.Error().Err(err).Msg("failed to process user photo")
				}

				messageData.Replies.Messages[index].FromID.ImageURL = userImageUrl
				messageData.Replies.Messages[index].FromID.Fullname = fmt.Sprintf("%s %s", replyData.FromID.FirstName, replyData.FromID.LastName)

				var replyImageUrl string

				if replyData.Media.Photo != nil {
					replyPhotoData, err := c.Replies.GetReplyPhoto(c.ctx, replyData)
					if err != nil {
						log.Error().Err(err).Msg("failed to get reply photo")
					}

					replyImageUrl, err = c.Photos.ProcessPhoto(c.ctx, replyPhotoData, fmt.Sprint(replyData.ID))
					if err != nil {
						log.Error().Err(err).Msg("failed to process reply photo")
					}
				}

				replyData.ImageURL = replyImageUrl
			}

			err = kafka.PushDataToQueue("messages", c.cfg.KafkaAddr, messageData)
			if err != nil {
				log.Error().Err(err).Msg("failed to push message into queue")
			}
		}

		time.Sleep(_saveTimeout)
	}
}

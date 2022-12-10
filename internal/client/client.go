package client

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gotd/td/tg"

	"github.com/VladPetriv/tg_scanner/internal/client/filter"
	"github.com/VladPetriv/tg_scanner/internal/client/group"
	"github.com/VladPetriv/tg_scanner/internal/client/message"
	"github.com/VladPetriv/tg_scanner/internal/client/photo"
	"github.com/VladPetriv/tg_scanner/internal/client/reply"
	"github.com/VladPetriv/tg_scanner/internal/client/user"
	"github.com/VladPetriv/tg_scanner/internal/controller"
	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/internal/store"
	"github.com/VladPetriv/tg_scanner/pkg/config"
	"github.com/VladPetriv/tg_scanner/pkg/file"
	"github.com/VladPetriv/tg_scanner/pkg/logger"
)

// Timeouts
var (
	_startTimeout         time.Duration = 20 * time.Second
	_historyTimeout       time.Duration = 30 * time.Minute
	_removeTimeout        time.Duration = 30 * time.Minute
	_saveTimeout          time.Duration = 15 * time.Minute
	_beetweenGroupTimeout time.Duration = 10 * time.Second
	_incomingTimeout      time.Duration = time.Minute
)

type appClient struct {
	ctx   context.Context
	store *store.Store
	queue controller.Controller
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

func New(ctx context.Context, store *store.Store, queue controller.Controller, api *tg.Client, log *logger.Logger, cfg *config.Config) *appClient {
	return &appClient{
		ctx:      ctx,
		store:    store,
		api:      api,
		log:      log,
		cfg:      cfg,
		queue:    queue,
		Groups:   group.New(log, api),
		Messages: message.New(log, api),
		Users:    user.New(log, api),
		Replies:  reply.New(log, api),
		Photos:   photo.New(log, store),
	}
}

func (c appClient) GetHistoryMessages(groups []model.TgGroup) {
	logger := c.log

	time.Sleep(_startTimeout)

	for {
		for _, group := range groups {
			logger.Info().Msgf("get - [%s]", group.Username)

			filePath := fmt.Sprintf("./data/%s.json", group.Username)

			groupPeer := &tg.InputPeerChannel{
				ChannelID:  group.ID,
				AccessHash: group.AccessHash,
			}

			tgMessages, err := c.Groups.GetMessagesFromGroupHistory(c.ctx, groupPeer)
			if err != nil {
				logger.Error().Err(err).Msg("get messages from group history")
			}

			modifiedTgMessages, ok := tgMessages.AsModified()
			if !ok {
				logger.Warn().Msg("receive unexpected messages type")
			}

			parsedMessages := c.Messages.ParseHistoryMessages(c.ctx, modifiedTgMessages, groupPeer)

			messages := make([]model.TgMessage, 0)

			for _, message := range parsedMessages {
				// check if message is question
				ok := filter.ProcessMessage(&message)
				if !ok {
					continue
				}

				message.PeerID = group

				user, err := c.Users.GetUser(c.ctx, message, groupPeer)
				if err != nil {
					logger.Error().Err(err).Msg("get user info for message")

					continue
				}

				message.FromID = *user

				tgReplies, err := c.Replies.GetReplies(c.ctx, message, groupPeer)
				if err != nil {
					logger.Error().Err(err).Msg("get replies for message")

					continue
				}

				parsedReplies := c.Replies.ParseTelegramReplies(c.ctx, tgReplies, groupPeer)

				// get user info for replies
				for index, reply := range parsedReplies {
					userInfo, err := c.Users.GetUser(c.ctx, reply, groupPeer)
					if err != nil {
						logger.Error().Err(err).Msg("get user info for reply")

						continue
					}

					parsedReplies[index].FromID = *userInfo
				}

				message.Replies.Count = len(parsedReplies)
				message.Replies.Messages = parsedReplies

				messages = append(messages, message)
			}

			messagesFromFile, err := c.Messages.GetMessagesFromFile(filePath)
			if err != nil {
				logger.Error().Err(err).Msg("get messages from the file")
			}

			messagesFromFile = append(messagesFromFile, messages...)

			filteredmessages := filter.RemoveDuplicatesFromMessages(messagesFromFile)

			c.Messages.WriteMessagesToFile(filteredmessages, filePath)

			time.Sleep(time.Second * 10)
		}

		time.Sleep(_historyTimeout)
	}
}

func (c appClient) GetIncomingMessages(user *tg.User, groups []model.TgGroup) {
	logger := c.log

	time.Sleep(_startTimeout)

	filePath := "./data/incoming.json"

	err := file.CreateFileForIncoming()
	if err != nil {
		logger.Error().Err(err).Msg("create base file for incoming messages")
	}

	for {
		processedMessages, err := c.Messages.ProcessIncomingMessages(c.ctx, user, groups)
		if err != nil {
			logger.Error().Err(err).Msg("process incoming messages")
		}

		messagesFromFile, err := c.Messages.GetMessagesFromFile(filePath)
		if err != nil {
			logger.Error().Err(err).Msg("get message from file")
		}

		for _, msg := range processedMessages {
			// check if message in question
			ok := filter.ProcessMessage(&msg)
			if !ok {
				continue
			}

			userInfo, err := c.Users.GetUser(c.ctx, msg, &tg.InputPeerChannel{
				ChannelID:  msg.PeerID.ID,
				AccessHash: msg.PeerID.AccessHash,
			})
			if err != nil {
				logger.Error().Err(err).Msg("get user info for message")

				continue
			}

			msg.FromID = *userInfo
		}

		messagesFromFile = append(messagesFromFile, processedMessages...)

		result := filter.RemoveDuplicatesFromMessages(messagesFromFile)

		c.Messages.WriteMessagesToFile(result, filePath)

		time.Sleep(_incomingTimeout)
	}
}

func (c appClient) PushToQueue() {
	logger := c.log

	for {
		messages, err := c.ProcessMessagesFromFiles("data")
		if err != nil {
			logger.Error().Err(err).Msg("get messages from files")
		}

		for _, messageData := range messages {
			messageValue, err := c.store.Cache.Get(c.ctx, c.store.Cache.GenerateKey(messageData))
			if err != nil {
				logger.Error().Err(err).Msg("get message key from cache")
			}

			if messageValue == "" {
				err := c.store.Cache.Set(c.ctx, c.store.Cache.GenerateKey(messageData), true)
				if err != nil {
					logger.Error().Err(err).Msg("set message into cache")
				}
			} else {
				logger.Info().Msg("message is exist")

				continue
			}

			groupPeer := &tg.InputPeerChannel{
				ChannelID:  messageData.PeerID.ID,
				AccessHash: messageData.PeerID.AccessHash,
			}

			replies, err := c.Replies.GetReplies(c.ctx, &messageData, groupPeer)
			if err != nil {
				logger.Error().Err(err).Msg("get replies for message")

				continue
			}

			processedReplies := c.Replies.ProcessReplies(c.ctx, replies, groupPeer)

			for _, reply := range processedReplies {
				userInfo, err := c.Users.GetUser(c.ctx, reply, groupPeer)
				if err != nil {
					logger.Error().Err(err).Msg("failed to get user info for reply")

					continue
				}

				reply.FromID = *userInfo
			}

			messageData.Replies.Count = len(processedReplies)
			messageData.Replies.Messages = processedReplies

			filter.RemoveDuplicatesFromReplies(&messageData.Replies)

			// if length of replies is 0 move to other message
			if len(messageData.Replies.Messages) == 0 {
				logger.Info().Msg("message have no replies")

				continue
			}

			userPhotoData, err := c.Users.GetUserPhoto(c.ctx, messageData.FromID)
			if err != nil {
				logger.Error().Err(err).Msg("get user photo")
			}

			userImageUrl, err := c.Photos.ProcessPhoto(c.ctx, userPhotoData, messageData.FromID.Username)
			if err != nil {
				logger.Error().Err(err).Msg("process user photo")
			}

			messageData.FromID.ImageURL = userImageUrl
			messageData.FromID.Fullname = fmt.Sprintf("%s %s", messageData.FromID.FirstName, messageData.FromID.LastName)

			var messageImageUrl string

			if ok, _ := c.Messages.CheckMessagePhotoStatus(c.ctx, &messageData); ok {
				messagePhotoData, err := c.Messages.GetMessagePhoto(c.ctx, messageData)
				if err != nil {
					logger.Error().Err(err).Msg("check message photo status")
				}

				messageImageUrl, err = c.Photos.ProcessPhoto(c.ctx, messagePhotoData, fmt.Sprint(messageData.ID))
				if err != nil {
					logger.Error().Err(err).Msg("process message photo")
				}
			}

			messageData.MessageURL = fmt.Sprintf("https://t.me/%s/%d", messageData.PeerID.Username, messageData.ID)
			messageData.ImageURL = messageImageUrl

			for index, replyData := range messageData.Replies.Messages {
				userPhotoData, err := c.Users.GetUserPhoto(c.ctx, replyData.FromID)
				if err != nil {
					logger.Error().Err(err).Msg("get user photo")
				}

				userImageUrl, err := c.Photos.ProcessPhoto(c.ctx, userPhotoData, replyData.FromID.Username)
				if err != nil {
					logger.Error().Err(err).Msg("process user photo")
				}

				messageData.Replies.Messages[index].FromID.ImageURL = userImageUrl
				messageData.Replies.Messages[index].FromID.Fullname = fmt.Sprintf("%s %s", replyData.FromID.FirstName, replyData.FromID.LastName)

				var replyImageUrl string

				if replyData.Media.Photo != nil {
					replyPhotoData, err := c.Replies.GetReplyPhoto(c.ctx, replyData)
					if err != nil {
						logger.Error().Err(err).Msg("get reply photo")
					}

					replyImageUrl, err = c.Photos.ProcessPhoto(c.ctx, replyPhotoData, fmt.Sprint(replyData.ID))
					if err != nil {
						logger.Error().Err(err).Msg("process reply photo")
					}
				}

				replyData.ImageURL = replyImageUrl
			}

			err = c.queue.PushDataToQueue("messages", messageData)
			if err != nil {
				logger.Error().Err(err).Msg("failed to push message into queue")
			}
		}

		time.Sleep(_saveTimeout)
	}
}

func (c appClient) ProcessMessagesFromFiles(path string) ([]model.TgMessage, error) {
	logger := c.log

	messages := make([]model.TgMessage, 0)

	directory, err := os.Open(path)
	if err != nil {
		logger.Error().Err(err).Msg("open directory")
		return nil, fmt.Errorf("open directory error: %w", err)
	}

	files, err := directory.ReadDir(0)
	if err != nil {
		logger.Error().Err(err).Msg("get all files from directory")
		return nil, fmt.Errorf("read directory error: %w", err)
	}

	for _, file := range files {
		filePath := fmt.Sprintf("./%s/%s", path, file.Name())

		messagesFromFile, err := c.Messages.GetMessagesFromFile(filePath)
		if err != nil {
			logger.Warn().Err(err).Msgf("get messages from file[%s]", file.Name())

			continue
		}

		/* err = os.Remove(pathToFile)
		if err != nil {
			return nil, fmt.Errorf("remove file error: %w", err)
		} */

		file, err := os.Create(filePath)
		if err != nil {
			logger.Error().Err(err).Msg("create file")
			return nil, fmt.Errorf("create file error: %w", err)
		}

		_, err = file.WriteString("[  ]")
		if err != nil {
			logger.Error().Err(err).Msg("write to file")
			return nil, fmt.Errorf("write to file error: %w", err)
		}

		messages = append(messages, messagesFromFile...)
	}

	result := filter.RemoveDuplicatesFromMessages(messages)

	return result, nil
}

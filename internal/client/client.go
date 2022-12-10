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

func (c appClient) GetIncomingMessages(tgUser tg.User, groups []model.TgGroup) {
	logger := c.log

	time.Sleep(_startTimeout)

	err := file.CreateFileForIncoming()
	if err != nil {
		logger.Error().Err(err).Msg("create base file for incoming messages")
	}

	for {
		parsedMessages, err := c.Messages.ParseIncomingMessages(c.ctx, tgUser, groups)
		if err != nil {
			logger.Error().Err(err).Msg("parse incoming messages from tg")
		}

		messages := make([]model.TgMessage, 0)

		for _, message := range parsedMessages {
			// check if message in question
			ok := filter.ProcessMessage(&message)
			if !ok {
				continue
			}

			//NOTE: add group info because incoming message don't have it
			for _, group := range groups {
				if message.PeerID.ChannelID == group.ID {
					message.PeerID = group
				}
			}

			user, err := c.Users.GetUser(c.ctx, message, &tg.InputPeerChannel{
				ChannelID:  message.PeerID.ID,
				AccessHash: message.PeerID.AccessHash,
			})
			if err != nil {
				logger.Error().Err(err).Msg("get user info for message")

				continue
			}

			message.FromID = *user

			messages = append(messages, message)
		}

		messagesFromFile, err := c.Messages.GetMessagesFromFile("./data/incoming.json")
		if err != nil {
			logger.Error().Err(err).Msg("get messages from file")
		}

		messagesFromFile = append(messagesFromFile, messages...)

		filteredmessages := filter.RemoveDuplicatesFromMessages(messagesFromFile)

		c.Messages.WriteMessagesToFile(filteredmessages, "./data/incoming.json")

		time.Sleep(_incomingTimeout)
	}
}

func (c appClient) ValidateAndPushGroupsToQueue(ctx context.Context) ([]model.TgGroup, error) {
	logger := c.log

	groups, err := c.Groups.GetGroups(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("get user groups")
		return nil, fmt.Errorf("failed to get groups: %w", err)
	}

	err = c.Groups.CreateFilesForGroups(groups)
	if err != nil {
		logger.Error().Err(err).Msg("create base files for groups")
		return nil, fmt.Errorf("failed to create base files for groups: %w", err)
	}

	for _, group := range groups {
		if group.ID == 0 {
			continue
		}

		isExist, err := c.store.Cache.Get(ctx, c.store.Cache.GenerateKey(group))
		if err != nil {
			logger.Error().Err(err).Msg("get value from cache by generated group key")
		}

		if isExist != "" {
			continue
		}

		err = c.store.Cache.Set(ctx, c.store.Cache.GenerateKey(group), true)
		if err != nil {
			logger.Error().Err(err).Msg("set value into cache with generated group key")
		}

		photo, err := c.Groups.GetGroupPhoto(ctx, &group)
		if err != nil {
			logger.Error().Err(err).Msgf("get [%s] photo data", group.Username)

			continue
		}

		imageURL, err := c.Photos.ProcessPhoto(ctx, photo, group.Username)
		if err != nil {
			logger.Error().Err(err).Msgf("process [%s] photo data", group.Username)
		}

		group.ImageURL = imageURL

		err = c.queue.PushDataToQueue("groups", group)
		if err != nil {
			logger.Error().Err(err).Msgf("push [%s] into queue", group.Username)
		}
	}

	return groups, nil
}

func (c appClient) PushMessagesToQueue() {
	logger := c.log

	for {
		messages, err := c.processMessagesFromFiles("data")
		if err != nil {
			logger.Error().Err(err).Msg("get messages from files")
		}

		for _, message := range messages {
			// check if message exist in cache
			messageCacheValue, err := c.store.Cache.Get(c.ctx, c.store.Cache.GenerateKey(message))
			if err != nil {
				logger.Error().Err(err).Msg("get message key from cache")
			}

			if messageCacheValue == "" {
				err := c.store.Cache.Set(c.ctx, c.store.Cache.GenerateKey(message), true)
				if err != nil {
					logger.Error().Err(err).Msg("set message into cache")
				}
			} else {
				logger.Info().Msg("message is exist")

				continue
			}

			groupPeer := &tg.InputPeerChannel{
				ChannelID:  message.PeerID.ID,
				AccessHash: message.PeerID.AccessHash,
			}

			replies, err := c.Replies.GetReplies(c.ctx, message, groupPeer)
			if err != nil {
				logger.Error().Err(err).Msg("get replies for message")

				continue
			}

			parsedReplies := c.Replies.ParseTelegramReplies(c.ctx, replies, groupPeer)

			for index, reply := range parsedReplies {
				userInfo, err := c.Users.GetUser(c.ctx, reply, groupPeer)
				if err != nil {
					logger.Error().Err(err).Msg("failed to get user info for reply")

					continue
				}

				parsedReplies[index].FromID = *userInfo
			}

			message.Replies.Count = len(parsedReplies)
			message.Replies.Messages = parsedReplies

			filter.RemoveDuplicatesFromReplies(&message.Replies)

			// if length of replies is 0 move to other message
			if len(message.Replies.Messages) == 0 {
				logger.Info().Msg("message have no replies")

				continue
			}

			c.processPhotosBeforePushToQueue(&message)

			message.FromID.Fullname = fmt.Sprintf(
				"%s %s",
				message.FromID.FirstName,
				message.FromID.LastName,
			)
			message.MessageURL = fmt.Sprintf(
				"https://t.me/%s/%d",
				message.PeerID.Username,
				message.ID,
			)

			for index, reply := range message.Replies.Messages {
				reply.FromID.Fullname = fmt.Sprintf(
					"%s %s",
					reply.FromID.FirstName,
					reply.FromID.LastName,
				)

				message.Replies.Messages[index].FromID.Fullname = fmt.Sprintf("%s %s", reply.FromID.FirstName, reply.FromID.LastName)
			}

			err = c.queue.PushDataToQueue("messages", message)
			if err != nil {
				logger.Error().Err(err).Msg("failed to push message into queue")
			}
		}

		time.Sleep(_saveTimeout)
	}
}

func (c appClient) processMessagesFromFiles(path string) ([]model.TgMessage, error) {
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

func (c appClient) processPhotosBeforePushToQueue(message *model.TgMessage) {
	logger := c.log

	photoData, err := c.Users.GetUserPhoto(c.ctx, message.FromID)
	if err != nil {
		logger.Error().Err(err).Msg("get user photo")
	}

	userImageUrl, err := c.Photos.ProcessPhoto(c.ctx, photoData, message.FromID.Username)
	if err != nil {
		logger.Error().Err(err).Msg("process user photo")
	}

	message.FromID.ImageURL = userImageUrl

	var messageImageUrl string

	if ok, _ := c.Messages.CheckMessagePhotoStatus(c.ctx, message); ok {
		photoData, err := c.Messages.GetMessagePhoto(c.ctx, *message)
		if err != nil {
			logger.Error().Err(err).Msg("check message photo status")
		}

		messageImageUrl, err = c.Photos.ProcessPhoto(c.ctx, photoData, fmt.Sprint(message.ID))
		if err != nil {
			logger.Error().Err(err).Msg("process message photo")
		}
	}

	message.ImageURL = messageImageUrl

	for index, reply := range message.Replies.Messages {
		photoData, err := c.Users.GetUserPhoto(c.ctx, reply.FromID)
		if err != nil {
			logger.Error().Err(err).Msg("get user photo")
		}

		userImageUrl, err := c.Photos.ProcessPhoto(c.ctx, photoData, reply.FromID.Username)
		if err != nil {
			logger.Error().Err(err).Msg("process user photo")
		}

		message.Replies.Messages[index].FromID.ImageURL = userImageUrl

		var replyImageUrl string

		if reply.Media.Photo != nil {
			photoData, err := c.Replies.GetReplyPhoto(c.ctx, reply)
			if err != nil {
				logger.Error().Err(err).Msg("get reply photo")
			}

			replyImageUrl, err = c.Photos.ProcessPhoto(c.ctx, photoData, fmt.Sprint(reply.ID))
			if err != nil {
				logger.Error().Err(err).Msg("process reply photo")
			}
		}

		reply.ImageURL = replyImageUrl
	}
}

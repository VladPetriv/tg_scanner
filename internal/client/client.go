package client

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gotd/td/tg"

	"github.com/VladPetriv/tg_scanner/config"
	"github.com/VladPetriv/tg_scanner/internal/client/filter"
	"github.com/VladPetriv/tg_scanner/internal/client/group"
	"github.com/VladPetriv/tg_scanner/internal/client/message"
	"github.com/VladPetriv/tg_scanner/internal/client/photo"
	"github.com/VladPetriv/tg_scanner/internal/client/reply"
	"github.com/VladPetriv/tg_scanner/internal/client/user"
	"github.com/VladPetriv/tg_scanner/internal/controller"
	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/internal/store"
	"github.com/VladPetriv/tg_scanner/pkg/file"
	"github.com/VladPetriv/tg_scanner/pkg/logger"
)

const (
	_startTimeout    = 20 * time.Second
	_historyTimeout  = 30 * time.Minute
	_saveTimeout     = 15 * time.Minute
	_betweenGroup    = 10 * time.Second
	_incomingTimeout = time.Minute
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

type AppClientOptions struct {
	Ctx   context.Context
	Store *store.Store
	Queue controller.Controller
	API   *tg.Client
	Log   *logger.Logger
	Cfg   *config.Config
}

var _ AppClient = (*appClient)(nil)

func New(options AppClientOptions) AppClient {
	return &appClient{
		ctx:      options.Ctx,
		store:    options.Store,
		api:      options.API,
		log:      options.Log,
		cfg:      options.Cfg,
		queue:    options.Queue,
		Groups:   group.New(options.Log, options.API),
		Messages: message.New(options.Log, options.API),
		Users:    user.New(options.Log, options.API),
		Replies:  reply.New(options.Log, options.API),
		Photos:   photo.New(options.Log, options.Store),
	}
}

func (c appClient) GetQuestionsFromGroupHistory(groups []model.Group) {
	logger := c.log

	time.Sleep(_startTimeout)

	for {
		for _, group := range groups {
			logger.Info().Msgf("get - [%s]", group.Username)

			filePath := fmt.Sprintf("./data/%s.json", group.Username)

			messages, err := c.Messages.GetHistoryMessagesFromGroup(c.ctx, group)
			if err != nil {
				logger.Error().Err(err).Msg("get messages from group")
			}

			questions := make([]model.Message, 0)

			for index, message := range messages {
				if !filter.IsQuestion(message.Message) {
					continue
				}

				messages[index].Message = filter.ReplaceUnexpectedSymbols(message.Message)

				questions = append(questions, message)
			}

			questions = c.addAdditionalDataForHistoryMessages(questions, group)

			messagesFromFile, err := c.Messages.GetMessagesFromFile(filePath)
			if err != nil {
				logger.Error().Err(err).Msg("get messages from the file")
			}

			messagesFromFile = append(messagesFromFile, questions...)

			filteredMessages := filter.RemoveDuplicatesFromMessages(messagesFromFile)

			c.Messages.WriteMessagesToFile(filteredMessages, filePath)

			time.Sleep(_betweenGroup)
		}

		time.Sleep(_historyTimeout)
	}
}

func (c appClient) addAdditionalDataForHistoryMessages(msgs []model.Message, group model.Group) []model.Message {
	logger := c.log

	messages := make([]model.Message, 0)

	for _, message := range msgs {
		message.PeerID = group

		userInfo, err := c.Users.GetUser(c.ctx, message, &group)
		if err != nil {
			logger.Error().Err(err).Msg("get user info for message")

			continue
		}

		message.FromID = *userInfo

		replies, err := c.Replies.GetReplies(c.ctx, message)
		if err != nil {
			logger.Error().Err(err).Msg("get replies for message")

			continue
		}

		for index, reply := range replies {
			userInfo, err := c.Users.GetUser(c.ctx, reply, &group)
			if err != nil {
				logger.Error().Err(err).Msg("get user info for reply")

				continue
			}

			replies[index].FromID = *userInfo
		}

		message.Replies.Count = len(replies)
		message.Replies.Messages = replies

		messages = append(messages, message)
	}

	return messages
}

func (c appClient) GetQuestionsFromIncomingMessages(tgUser tg.User, groups []model.Group) {
	logger := c.log

	time.Sleep(_startTimeout)

	err := file.CreateFileForIncoming()
	if err != nil {
		logger.Error().Err(err).Msg("create base file for incoming messages")
	}

	for {
		logger.Info().Msg("get - [incoming messages]")

		messages, err := c.Messages.GetIncomingMessagesFromUserGroups(c.ctx, tgUser, groups)
		if err != nil {
			logger.Error().Err(err).Msg("get incoming messages from user group")
		}

		questions := make([]model.Message, 0)

		for _, message := range messages {
			// We won't save messages that are reply to other messages
			if message.ReplyTo.ReplyToMsgID != 0 {
				continue
			}

			if !filter.IsQuestion(message.Message) {
				continue
			}

			message.Message = filter.ReplaceUnexpectedSymbols(message.Message)

			questions = append(questions, message)
		}

		questions = c.addAdditionalDataForIncomingMessages(questions, groups)

		messagesFromFile, err := c.Messages.GetMessagesFromFile("./data/incoming.json")
		if err != nil {
			logger.Error().Err(err).Msg("get messages from file")
		}

		messagesFromFile = append(messagesFromFile, questions...)

		filteredMessages := filter.RemoveDuplicatesFromMessages(messagesFromFile)

		c.Messages.WriteMessagesToFile(filteredMessages, "./data/incoming.json")

		time.Sleep(_incomingTimeout)
	}
}

func (c appClient) addAdditionalDataForIncomingMessages(msgs []model.Message, groups []model.Group) []model.Message {
	logger := c.log

	messages := make([]model.Message, 0)

	for _, message := range msgs {
		// Add group info because incoming message don't have it
		for _, group := range groups {
			if message.PeerID.ChannelID == group.ID {
				message.PeerID = group
			}
		}

		user, err := c.Users.GetUser(c.ctx, message, &message.PeerID)
		if err != nil {
			logger.Error().Err(err).Msg("get user info for message")

			continue
		}

		message.FromID = *user

		messages = append(messages, message)
	}

	return messages
}

func (c appClient) ValidateAndPushGroupsToQueue(ctx context.Context) ([]model.Group, error) {
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

		groupHash, err := group.GetHash()
		if err != nil {
			logger.Error().Err(err).Msg("generate hash for group")

			continue
		}

		isExist, err := c.store.Cache.Get(ctx, groupHash)
		if err != nil {
			logger.Error().Err(err).Msg("get value from cache by generated group key")
		}

		if isExist != "" {
			continue
		}

		err = c.store.Cache.Set(ctx, groupHash, true)
		if err != nil {
			logger.Error().Err(err).Msg("set value into cache with generated group key")
		}

		photo, err := c.Groups.GetGroupPhoto(ctx, group)
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

func (c appClient) PushMessagesToQueue() { //nolint:gocognit
	logger := c.log

	for {
		messages, err := c.processMessagesFromFiles("data")
		if err != nil {
			logger.Error().Err(err).Msg("get messages from files")
		}

		for _, message := range messages {
			messageHash, err := message.GetHash()
			if err != nil {
				logger.Error().Err(err).Msg("generate hash for message")

				continue
			}

			// check if message exist in cache
			messageCacheValue, err := c.store.Cache.Get(c.ctx, messageHash)
			if err != nil {
				logger.Error().Err(err).Msg("get message key from cache")
			}

			if messageCacheValue == "" {
				err := c.store.Cache.Set(c.ctx, messageHash, true)
				if err != nil {
					logger.Error().Err(err).Msg("set message into cache")
				}
			} else {
				logger.Info().Msg("message is exist")

				continue
			}

			replies, err := c.Replies.GetReplies(c.ctx, message)
			if err != nil {
				logger.Error().Err(err).Msg("get replies for message")

				continue
			}

			for index, reply := range replies {
				userInfo, err := c.Users.GetUser(c.ctx, reply, &message.PeerID)
				if err != nil {
					logger.Error().Err(err).Msg("failed to get user info for reply")

					continue
				}

				replies[index].FromID = *userInfo
			}

			message.Replies.Count = len(replies)
			message.Replies.Messages = replies

			filter.RemoveDuplicatesFromReplies(&message.Replies)

			// if length of replies is 0 move to other message
			if len(message.Replies.Messages) == 0 {
				logger.Info().Msg("message have no replies")

				continue
			}

			c.processPhotosBeforePushToQueue(&message) //nolint:gosec// ...

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

				message.Replies.Messages[index].FromID.Fullname = fmt.Sprintf(
					"%s %s", reply.FromID.FirstName, reply.FromID.LastName,
				)
			}

			err = c.queue.PushDataToQueue("messages", message)
			if err != nil {
				logger.Error().Err(err).Msg("failed to push message into queue")
			}
		}

		time.Sleep(_saveTimeout)
	}
}

func (c appClient) processMessagesFromFiles(path string) ([]model.Message, error) {
	logger := c.log

	messages := make([]model.Message, 0)

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

func (c appClient) processPhotosBeforePushToQueue(message *model.Message) {
	logger := c.log

	userPhotoData, err := c.Users.GetUserPhoto(c.ctx, message.FromID)
	if err != nil {
		logger.Error().Err(err).Msg("get user photo")
	}

	userImageURL, err := c.Photos.ProcessPhoto(c.ctx, userPhotoData, message.FromID.Username)
	if err != nil {
		logger.Error().Err(err).Msg("process user photo")
	}

	message.FromID.ImageURL = userImageURL

	var messageImageURL string

	if ok, _ := c.Messages.CheckMessagePhotoStatus(c.ctx, message); ok {
		messagePhotoData, err := c.Messages.GetMessagePhoto(c.ctx, *message)
		if err != nil {
			logger.Error().Err(err).Msg("check message photo status")
		}

		messageImageURL, err = c.Photos.ProcessPhoto(c.ctx, messagePhotoData, fmt.Sprint(message.ID))
		if err != nil {
			logger.Error().Err(err).Msg("process message photo")
		}
	}

	message.ImageURL = messageImageURL

	for index, reply := range message.Replies.Messages {
		replyUserPhotoData, err := c.Users.GetUserPhoto(c.ctx, reply.FromID)
		if err != nil {
			logger.Error().Err(err).Msg("get user photo")
		}

		replyUserImageURL, err := c.Photos.ProcessPhoto(c.ctx, replyUserPhotoData, reply.FromID.Username)
		if err != nil {
			logger.Error().Err(err).Msg("process user photo")
		}

		message.Replies.Messages[index].FromID.ImageURL = replyUserImageURL

		var replyImageURL string

		if reply.Media.Photo != nil {
			replyPhotoData, err := c.Replies.GetReplyPhoto(c.ctx, reply)
			if err != nil {
				logger.Error().Err(err).Msg("get reply photo")
			}

			replyImageURL, err = c.Photos.ProcessPhoto(c.ctx, replyPhotoData, fmt.Sprint(reply.ID))
			if err != nil {
				logger.Error().Err(err).Msg("process reply photo")
			}
		}

		reply.ImageURL = replyImageURL
	}
}

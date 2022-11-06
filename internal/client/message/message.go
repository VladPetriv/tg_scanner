package message

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gotd/td/tg"

	"github.com/VladPetriv/tg_scanner/internal/client/photo"
	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/pkg/logger"
)

type tgMessage struct {
	log *logger.Logger
	api *tg.Client
}

var _ Message = (*tgMessage)(nil)

func New(log *logger.Logger, api *tg.Client) *tgMessage {
	return &tgMessage{
		log: log,
		api: api,
	}
}

func (m tgMessage) ProcessHistoryMessages(ctx context.Context, data tg.ModifiedMessagesMessages, groupPeer *tg.InputPeerChannel) []model.TgMessage {
	logger := m.log

	processedMessages := make([]model.TgMessage, 0)
	messages := data.GetMessages()

	for _, message := range messages {
		msg := model.TgMessage{}

		encodedData, err := json.Marshal(message)
		if err != nil {
			logger.Warn().Err(err).Msg("marshal message data")

			continue
		}

		err = json.Unmarshal(encodedData, &msg)
		if err != nil {
			logger.Warn().Err(err).Msg("unmarshal message data")

			continue
		}

		processedMessages = append(processedMessages, msg)
	}

	return processedMessages
}

func (m tgMessage) ProcessIncomingMessages(ctx context.Context, tgUser *tg.User, groups []model.TgGroup) ([]model.TgMessage, error) {
	logger := m.log
	processedMessages := make([]model.TgMessage, 0)

	data, err := m.api.MessagesGetDialogs(ctx, &tg.MessagesGetDialogsRequest{
		OffsetPeer: &tg.InputPeerUser{
			UserID:     tgUser.ID,
			AccessHash: tgUser.AccessHash,
		},
	})
	if err != nil {
		logger.Error().Err(err).Msg("get incoming messages")
		return nil, fmt.Errorf("get incoming messages error: %w", err)
	}

	modifiedMessages, _ := data.AsModified()

	for _, message := range modifiedMessages.GetMessages() {
		msg := model.TgMessage{}

		encodedData, err := json.Marshal(message)
		if err != nil {
			logger.Warn().Err(err).Msg("marshal message data")

			continue
		}

		err = json.Unmarshal(encodedData, &msg)
		if err != nil {
			logger.Warn().Err(err).Msg("unmarshal message data")

			continue
		}

		// add group info because incoming message don't have it
		for _, groupData := range groups {
			if msg.PeerID.ChannelID == groupData.ID {
				msg.PeerID = groupData
			}
		}

		processedMessages = append(processedMessages, msg)
	}

	return processedMessages, nil
}

func (m tgMessage) GetMessagePhoto(ctx context.Context, message model.TgMessage) (tg.UploadFileClass, error) {
	logger := m.log

	length := len(message.Media.Photo.Sizes) - 1

	data, err := m.api.UploadGetFile(ctx, &tg.UploadGetFileRequest{
		Location: &tg.InputPhotoFileLocation{
			ID:            message.Media.Photo.ID,
			AccessHash:    message.Media.Photo.AccessHash,
			FileReference: message.Media.Photo.FileReference,
			ThumbSize:     message.Media.Photo.Sizes[length].GetType(),
		},
		Offset: 0,
		Limit:  photo.Size,
	})
	if err != nil {
		logger.Error().Err(err).Msg("get message photo")
		return nil, fmt.Errorf("get message photo error: %w", err)
	}

	return data, nil
}

// TODO: refactor it
func (m tgMessage) CheckMessagePhotoStatus(ctx context.Context, message *model.TgMessage) (bool, error) {
	logger := m.log

	request := &tg.ChannelsGetMessagesRequest{
		Channel: &tg.InputChannel{
			ChannelID:  message.PeerID.ID,
			AccessHash: message.PeerID.AccessHash,
		},
		ID: []tg.InputMessageClass{&tg.InputMessageID{ID: message.ID}},
	}

	data, err := m.api.ChannelsGetMessages(ctx, request)
	if err != nil {
		logger.Error().Err(err).Msg("get message by id")
		return false, fmt.Errorf("get message by id error: %w", err)
	}

	messages, _ := data.(*tg.MessagesChannelMessages)

	for _, m := range messages.GetMessages() {
		message, ok := m.(*tg.Message)
		if !ok {
			logger.Warn().Bool("isMessage", ok).Msg("receive unexpected message type")

			continue
		}

		if message.Media != nil {
			media, ok := message.Media.(*tg.MessageMediaPhoto)
			if !ok {
				logger.Warn().Bool("isMedia", ok).Msg("receive unexpected media type")

				continue
			}

			photo, ok := media.GetPhoto()
			if !ok {
				logger.Warn().Bool("isPhoto", ok).Msg("receive unexpected photo type")

				continue
			}

			if photo != nil {
				return true, nil
			}

			return false, nil
		}
	}

	return false, nil
}

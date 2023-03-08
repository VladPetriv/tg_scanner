package message

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

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

func New(log *logger.Logger, api *tg.Client) Message {
	return &tgMessage{
		log: log,
		api: api,
	}
}

func (m tgMessage) GetHistoryMessagesFromGroup(ctx context.Context, group model.Group) ([]model.Message, error) {
	logger := m.log

	data, err := m.api.MessagesGetHistory(ctx, &tg.MessagesGetHistoryRequest{
		Peer: &tg.InputPeerChannel{
			ChannelID:  group.ID,
			AccessHash: group.AccessHash,
		},
	})
	if err != nil {
		logger.Error().Err(err).Msg("get messages from group history")
		return nil, fmt.Errorf("get messages from group history: %w", err)
	}

	modifiedData, ok := data.AsModified()
	if !ok {
		logger.Info().Msg("received unexpected type of messages")
		return nil, nil
	}

	tgMessages := modifiedData.GetMessages()

	messages := m.parseHistoryMessages(tgMessages)

	return messages, nil
}

func (m tgMessage) parseHistoryMessages(tgMessages []tg.MessageClass) []model.Message {
	logger := m.log

	messages := make([]model.Message, 0)

	for _, message := range tgMessages {
		var msg model.Message

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

		messages = append(messages, msg)
	}

	return messages
}

func (m tgMessage) GetIncomingMessagesFromUserGroups(ctx context.Context, tgUser tg.User, groups []model.Group) ([]model.Message, error) { //nolint: lll
	logger := m.log

	data, err := m.api.MessagesGetDialogs(ctx, &tg.MessagesGetDialogsRequest{
		OffsetPeer: &tg.InputPeerUser{
			UserID:     tgUser.ID,
			AccessHash: tgUser.AccessHash,
		},
	})
	if err != nil {
		logger.Error().Err(err).Msg("get incoming messages")
		return nil, fmt.Errorf("get incoming messages: %w", err)
	}

	modifiedData, ok := data.AsModified()
	if !ok {
		logger.Info().Msg("received unexpected type of messages")
	}

	tgMessages := modifiedData.GetMessages()

	messages := m.parseIncomingMessages(tgMessages)

	return messages, nil
}

func (m tgMessage) parseIncomingMessages(tgMessages []tg.MessageClass) []model.Message {
	logger := m.log

	messages := make([]model.Message, 0)

	for _, message := range tgMessages {
		var msg model.Message

		encodedData, err := json.Marshal(message)
		if err != nil {
			logger.Error().Err(err).Msg("marshal message data")

			continue
		}

		err = json.Unmarshal(encodedData, &msg)
		if err != nil {
			logger.Error().Err(err).Msg("unmarshal message data")

			continue
		}

		messages = append(messages, msg)
	}

	return messages
}

func (m tgMessage) GetMessagePhoto(ctx context.Context, message model.Message) (tg.UploadFileClass, error) {
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

func (m tgMessage) CheckMessagePhotoStatus(ctx context.Context, message *model.Message) (bool, error) {
	logger := m.log

	data, err := m.api.ChannelsGetMessages(ctx,
		&tg.ChannelsGetMessagesRequest{
			Channel: &tg.InputChannel{
				ChannelID:  message.PeerID.ID,
				AccessHash: message.PeerID.AccessHash,
			},
			ID: []tg.InputMessageClass{&tg.InputMessageID{ID: message.ID}},
		})
	if err != nil {
		logger.Error().Err(err).Msg("get message by id")
		return false, fmt.Errorf("get message by id error: %w", err)
	}

	messages, ok := data.(*tg.MessagesChannelMessages)

	if ok {
		//NOTE: we should check only first message
		message, isMessage := messages.GetMessages()[0].(*tg.Message)
		if !isMessage {
			logger.Warn().Bool("isMessage", isMessage).Msg("receive unexpected message type")

			return false, nil
		}

		if message.Media != nil {
			media, isMedia := message.Media.(*tg.MessageMediaPhoto)
			if !isMedia {
				logger.Warn().Bool("isMedia", isMedia).Msg("receive unexpected media type")

				return false, nil
			}

			_, isPhoto := media.GetPhoto()

			return isPhoto, nil
		}
	}

	return false, nil
}

func (m tgMessage) WriteMessagesToFile(messages []model.Message, fileName string) {
	logger := m.log

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		logger.Error().Err(err).Msg("open file error")
	}

	encodedMessages, err := json.Marshal(messages)
	if err != nil {
		logger.Error().Err(err).Msg("marshal messages error")
	}

	_, err = file.Write(encodedMessages)
	if err != nil {
		logger.Error().Err(err).Msg("write messages to file error")
	}
}

func (m tgMessage) GetMessagesFromFile(filePath string) ([]model.Message, error) {
	logger := m.log

	messages := make([]model.Message, 0)

	data, err := os.ReadFile(filePath)
	if err != nil {
		logger.Error().Err(err).Msg("read file")
		return nil, fmt.Errorf("read file error: %w", err)
	}

	err = json.Unmarshal(data, &messages)
	if err != nil {
		logger.Error().Err(err).Msg("unmarshal messages")
		return nil, fmt.Errorf("unmarshal file data error: %w", err)
	}

	_, err = os.Create(filePath)
	if err != nil {
		logger.Error().Err(err).Msg("create file")
		return nil, fmt.Errorf("create file error: %w", err)
	}

	return messages, nil
}

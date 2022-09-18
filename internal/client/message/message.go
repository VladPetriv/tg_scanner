package message

import (
	"context"
	"encoding/json"

	"github.com/gotd/td/tg"

	"github.com/VladPetriv/tg_scanner/internal/client/photo"
	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/pkg/errors"
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
	processedMessages := make([]model.TgMessage, 0)
	messages := data.GetMessages()

	for _, message := range messages {
		msg := model.TgMessage{}

		encodedData, err := json.Marshal(message)
		if err != nil {
			m.log.Warn().Err(err)

			continue
		}

		err = json.Unmarshal(encodedData, &msg)
		if err != nil {
			m.log.Warn().Err(err)

			continue
		}

		processedMessages = append(processedMessages, msg)
	}

	return processedMessages
}

func (m tgMessage) ProcessIncomingMessages(ctx context.Context, tgUser *tg.User, groups []model.TgGroup) ([]model.TgMessage, error) {
	processedMessages := make([]model.TgMessage, 0)

	data, err := m.api.MessagesGetDialogs(ctx, &tg.MessagesGetDialogsRequest{
		OffsetPeer: &tg.InputPeerUser{
			UserID:     tgUser.ID,
			AccessHash: tgUser.AccessHash,
		},
	})
	if err != nil {
		m.log.Error().Err(err)

		return nil, &errors.GetError{Name: "incoming messages", ErrorValue: err}
	}

	modifiedMessages, _ := data.AsModified()

	for _, message := range modifiedMessages.GetMessages() {
		msg := model.TgMessage{}

		encodedData, err := json.Marshal(message)
		if err != nil {
			m.log.Warn().Err(err)

			continue
		}

		err = json.Unmarshal(encodedData, &msg)
		if err != nil {
			m.log.Warn().Err(err)

			continue
		}

		// add group info because incoming message don't have it
		for _, channel := range groups {
			if msg.PeerID.ChannelID == channel.ID {
				msg.PeerID = channel
			}
		}

		processedMessages = append(processedMessages, msg)
	}

	return processedMessages, nil
}

func (m tgMessage) GetMessagePhoto(ctx context.Context, message model.TgMessage) (tg.UploadFileClass, error) {
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
		m.log.Error().Err(err)

		return nil, &errors.GetError{Name: "message photo", ErrorValue: err}
	}

	return data, nil
}

// TODO: refactor it
func (m tgMessage) CheckMessagePhotoStatus(ctx context.Context, message *model.TgMessage) (bool, error) {
	request := &tg.ChannelsGetMessagesRequest{
		Channel: &tg.InputChannel{
			ChannelID:  message.PeerID.ID,
			AccessHash: message.PeerID.AccessHash,
		},
		ID: []tg.InputMessageClass{&tg.InputMessageID{ID: message.ID}},
	}

	data, err := m.api.ChannelsGetMessages(ctx, request)
	if err != nil {
		return false, &errors.GetError{Name: "messages by channel peer", ErrorValue: err}
	}

	messages, _ := data.(*tg.MessagesChannelMessages)

	for _, m := range messages.GetMessages() {
		message, ok := m.(*tg.Message)
		if !ok {
			continue
		}

		if message.Media != nil {
			media, ok := message.Media.(*tg.MessageMediaPhoto)
			if !ok {
				continue
			}

			photo, ok := media.GetPhoto()
			if !ok {
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

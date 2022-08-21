package message

import (
	"context"
	"encoding/json"

	"github.com/gotd/td/tg"

	"github.com/VladPetriv/tg_scanner/internal/client/photo"
	"github.com/VladPetriv/tg_scanner/internal/client/replie"
	"github.com/VladPetriv/tg_scanner/internal/client/user"
	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/pkg/errors"
)

func GetMessagesFromTelegram(ctx context.Context, data tg.ModifiedMessagesMessages, channelPeer *tg.InputPeerChannel, api *tg.Client) []model.TgMessage { // nolint
	var msg model.TgMessage

	messages := make([]model.TgMessage, 0)
	messagesFromTg := data.GetMessages()

	for _, message := range messagesFromTg {
		encodedData, err := json.Marshal(message)
		if err != nil {
			continue
		}

		err = json.Unmarshal(encodedData, &msg)
		if err != nil {
			continue
		}

		replies, err := replie.GetReplies(ctx, &msg, channelPeer, api)
		if err != nil {
			continue
		}

		msg.Replies.Messages = replie.ProcessRepliesMessage(ctx, replies, channelPeer, api)
		msg.Replies.Count = len(msg.Replies.Messages)

		messages = append(messages, msg)
	}

	return messages
}

func GetIncomingMessages(ctx context.Context, tgUser *tg.User, channels []model.TgChannel, api *tg.Client) ([]model.TgMessage, error) {
	messages := make([]model.TgMessage, 0)

	data, err := api.MessagesGetDialogs(ctx, &tg.MessagesGetDialogsRequest{
		OffsetPeer: &tg.InputPeerUser{
			UserID:     tgUser.ID,
			AccessHash: tgUser.AccessHash,
		},
	})
	if err != nil {
		return nil, &errors.GettingError{Name: "incoming messages", ErrorValue: err}
	}

	modifiedData, _ := data.AsModified()
	for _, msg := range modifiedData.GetMessages() {
		message := model.TgMessage{}

		encodedData, err := json.Marshal(msg)
		if err != nil {
			continue
		}

		err = json.Unmarshal(encodedData, &message)
		if err != nil {
			continue
		}

		for _, channel := range channels {
			if message.PeerID.ChannelID == channel.ID {
				message.PeerID = channel
			}
		}

		userInfo, err := user.GetUserInfo(ctx, message.FromID.UserID, message.ID, &tg.InputPeerChannel{
			ChannelID:  message.PeerID.ID,
			AccessHash: message.PeerID.AccessHash,
		}, api)
		if err != nil {
			continue
		}

		message.FromID = *userInfo

		messages = append(messages, message)
	}

	return messages, nil
}

func GetMessagePhoto(ctx context.Context, msg model.TgMessage, api *tg.Client) (tg.UploadFileClass, error) {
	length := len(msg.Media.Photo.Sizes) - 1

	data, err := api.UploadGetFile(ctx, &tg.UploadGetFileRequest{
		Location: &tg.InputPhotoFileLocation{
			ID:            msg.Media.Photo.ID,
			AccessHash:    msg.Media.Photo.AccessHash,
			FileReference: msg.Media.Photo.FileReference,
			ThumbSize:     msg.Media.Photo.Sizes[length].GetType(),
		},
		Offset: 0,
		Limit:  photo.Size,
	})
	if err != nil {
		return nil, &errors.GettingError{Name: "message photo", ErrorValue: err}
	}

	return data, nil
}

func CheckMessagePhotoStatus(ctx context.Context, msg *model.TgMessage, api *tg.Client) (bool, error) {
	request := &tg.ChannelsGetMessagesRequest{
		Channel: &tg.InputChannel{
			ChannelID:  msg.PeerID.ID,
			AccessHash: msg.PeerID.AccessHash,
		},
		ID: []tg.InputMessageClass{&tg.InputMessageID{ID: msg.ID}},
	}

	data, err := api.ChannelsGetMessages(ctx, request)
	if err != nil {
		return false, &errors.GettingError{Name: "messages by channel peer", ErrorValue: err}
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

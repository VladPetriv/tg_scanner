package message

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gotd/td/tg"

	"github.com/VladPetriv/tg_scanner/internal/client/user"
	"github.com/VladPetriv/tg_scanner/internal/file"
	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/pkg/utils"
)

var messageImageSize int = 1024 * 1024

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

		replies, err := GetReplies(ctx, &msg, channelPeer, api)
		if err != nil {
			continue
		}

		msg.Replies.Messages = ProcessRepliesMessage(ctx, replies, channelPeer, api)
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
		return nil, &utils.GettingError{Name: "incoming messages", ErrorValue: err}
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
			if message.PeerID.ChannelID == channel.ID && message.PeerID.Username == channel.Username {
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

func GetReplies(ctx context.Context, message *model.TgMessage, channelPeer *tg.InputPeerChannel, api *tg.Client) (tg.MessagesMessagesClass, error) { // nolint
	replies, err := api.MessagesGetReplies(ctx, &tg.MessagesGetRepliesRequest{
		Peer:  channelPeer,
		MsgID: message.ID,
	})
	if err != nil {
		return nil, &utils.GettingError{Name: "replies", ErrorValue: err}
	}

	return replies, nil
}

func GetRepliesForMessageBeforeSave(ctx context.Context, message *model.TgMessage, api *tg.Client) error {
	channelPeer := &tg.InputPeerChannel{
		ChannelID:  message.PeerID.ID,
		AccessHash: message.PeerID.AccessHash,
	}

	replies, err := GetReplies(ctx, message, channelPeer, api)
	if err != nil {
		return err
	}

	messageReplie := ProcessRepliesMessage(ctx, replies, channelPeer, api)

	message.Replies.Messages = append(message.Replies.Messages, messageReplie...)

	time.Sleep(time.Second * 3)

	return nil
}

func ProcessRepliesMessage(ctx context.Context, replies tg.MessagesMessagesClass, cPeer *tg.InputPeerChannel, api *tg.Client) []model.TgRepliesMessage {
	repliesMessages := make([]model.TgRepliesMessage, 0)

	data, _ := replies.AsModified()
	for _, replie := range data.GetMessages() {
		replieMessage := model.TgRepliesMessage{}

		encodedData, err := json.Marshal(replie)
		if err != nil {
			continue
		}

		err = json.Unmarshal(encodedData, &replieMessage)
		if err != nil {
			continue
		}

		userInfo, err := user.GetUserInfo(ctx, replieMessage.FromID.UserID, replieMessage.ID, cPeer, api)
		if err != nil {
			continue
		}

		replieMessage.FromID = *userInfo

		repliesMessages = append(repliesMessages, replieMessage)
	}

	return repliesMessages
}

func GetMessagePhoto(ctx context.Context, msg *model.TgMessage, api *tg.Client) (tg.UploadFileClass, error) {
	length := len(msg.Media.Photo.Sizes) - 1

	data, err := api.UploadGetFile(ctx, &tg.UploadGetFileRequest{
		Location: &tg.InputPhotoFileLocation{
			ID:            msg.Media.Photo.ID,
			AccessHash:    msg.Media.Photo.AccessHash,
			FileReference: msg.Media.Photo.FileReference,
			ThumbSize:     msg.Media.Photo.Sizes[length].GetType(),
		},
		Offset: 0,
		Limit:  messageImageSize,
	})
	if err != nil {
		return nil, &utils.GettingError{Name: "message photo", ErrorValue: err}
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
		return false, &utils.GettingError{Name: "messages by channel peer", ErrorValue: err}
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

func ProcessMessagePhoto(ctx context.Context, msg *model.TgMessage, api *tg.Client) (string, error) {
	status, err := CheckMessagePhotoStatus(ctx, msg, api)
	if err != nil {
		return "", err
	}

	if !status {
		return "", &utils.NotFoundError{Name: "photo in message"}
	}

	messagePhotoData, err := GetMessagePhoto(ctx, msg, api)
	if err != nil {
		return "", err
	}

	messageImage, err := file.DecodePhoto(messagePhotoData)
	if err != nil {
		return "", fmt.Errorf("decode message photo error: %w", err)
	}

	filename, err := file.CreatePhoto(messageImage, fmt.Sprint(msg.ID))
	if err != nil {
		return "", err
	}

	return filename, err
}

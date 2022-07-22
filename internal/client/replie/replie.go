package replie

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gotd/td/tg"

	"github.com/VladPetriv/tg_scanner/internal/client/photo"
	"github.com/VladPetriv/tg_scanner/internal/client/user"
	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/pkg/utils"
)

func GetReplies(ctx context.Context, message *model.TgMessage, channelPeer *tg.InputPeerChannel, api *tg.Client) (tg.MessagesMessagesClass, error) {
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

func GetRepliePhoto(ctx context.Context, replie model.TgRepliesMessage, api *tg.Client) (tg.UploadFileClass, error) {
	length := len(replie.Media.Photo.Sizes) - 1

	data, err := api.UploadGetFile(ctx, &tg.UploadGetFileRequest{
		Location: &tg.InputPhotoFileLocation{
			ID:            replie.Media.Photo.ID,
			AccessHash:    replie.Media.Photo.AccessHash,
			FileReference: replie.Media.Photo.FileReference,
			ThumbSize:     replie.Media.Photo.Sizes[length].GetType(),
		},
		Offset: 0,
		Limit:  photo.Size,
	})
	if err != nil {
		return nil, &utils.GettingError{Name: "replie photo", ErrorValue: err}
	}

	return data, nil
}

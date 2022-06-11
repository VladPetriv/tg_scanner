package message

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"math/big"
	"time"

	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/internal/user"
	"github.com/VladPetriv/tg_scanner/pkg/utils"
	"github.com/gotd/td/tg"
)

func GetMessagesFromTelegram(ctx context.Context, data tg.ModifiedMessagesMessages, channelPeer *tg.InputPeerChannel, api *tg.Client) []model.TgMessage { // nolint
	var msg model.TgMessage

	result := make([]model.TgMessage, 0)
	messages := data.GetMessages()

	for _, message := range messages {
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

		repliesMessages := ProcessRepliesMessage(ctx, replies, channelPeer, api)
		msg.Replies.Count = len(repliesMessages)
		msg.Replies.Messages = repliesMessages

		result = append(result, msg)
	}

	return result
}

func GetReplies(ctx context.Context, message *model.TgMessage, channelPeer *tg.InputPeerChannel, api *tg.Client) (tg.MessagesMessagesClass, error) { // nolint
	bInt := big.NewInt(10000) // nolint

	value, _ := rand.Int(rand.Reader, bInt)

	replies, err := api.MessagesGetReplies(ctx, &tg.MessagesGetRepliesRequest{ // nolint
		Peer:  channelPeer,
		MsgID: message.ID,
		Hash:  value.Int64(),
	})
	if err != nil {
		return nil, &utils.GettingError{Name: "replies", ErrorValue: err}
	}

	return replies, nil
}

func ProcessRepliesMessage(ctx context.Context, replies tg.MessagesMessagesClass, cPeer *tg.InputPeerChannel, api *tg.Client) []model.TgRepliesMessage {
	repliesMessages := make([]model.TgRepliesMessage, 0)

	var replieMessage model.TgRepliesMessage

	data, _ := replies.AsModified()
	for _, replie := range data.GetMessages() {
		encodedData, err := json.Marshal(replie)
		if err != nil {
			continue
		}

		err = json.Unmarshal(encodedData, &replieMessage)
		if err != nil {
			continue
		}

		u, err := user.GetUserInfo(ctx, replieMessage.FromID.UserID, replieMessage.ID, cPeer, api)
		if err != nil {
			continue
		}

		replieMessage.FromID = *u
		repliesMessages = append(repliesMessages, replieMessage)
	}

	return repliesMessages
}

func GetIncomingMessages(ctx context.Context, tg_user *tg.User, channels []model.TgChannel, api *tg.Client) ([]model.TgMessage, error) {
	msgs := make([]model.TgMessage, 0)

	var msg model.TgMessage

	data, err := api.MessagesGetDialogs(ctx, &tg.MessagesGetDialogsRequest{ // nolint
		OffsetPeer: &tg.InputPeerUser{
			UserID:     tg_user.ID,
			AccessHash: tg_user.AccessHash,
		},
	})
	if err != nil {
		return nil, &utils.GettingError{Name: "incoming messages", ErrorValue: err}
	}

	modifiedData, _ := data.AsModified()
	for _, m := range modifiedData.GetMessages() {
		encodedData, err := json.Marshal(m)
		if err != nil {
			continue
		}

		err = json.Unmarshal(encodedData, &msg)
		if err != nil {
			continue
		}

		// Getting channel info for replie
		for _, channel := range channels {
			if msg.PeerID.ChannelID == channel.ID {
				msg.PeerID = channel
			}
		}

		// Getting user info for replie
		u, err := user.GetUserInfo(ctx, msg.FromID.UserID, msg.ID, &tg.InputPeerChannel{
			ChannelID:  int64(msg.PeerID.ID),
			AccessHash: int64(msg.PeerID.AccessHash),
		}, api)
		if err != nil {
			continue
		}

		msg.FromID = *u

		msgs = append(msgs, msg)
	}

	return msgs, nil
}

func GetRepliesForMessageBeforeSave(ctx context.Context, message *model.TgMessage, api *tg.Client) error {
	cPeer := &tg.InputPeerChannel{
		ChannelID:  int64(message.PeerID.ID),
		AccessHash: int64(message.PeerID.AccessHash),
	}

	replies, err := GetReplies(ctx, message, cPeer, api)
	if err != nil {
		return err
	}

	messageReplie := ProcessRepliesMessage(ctx, replies, cPeer, api)

	message.Replies.Messages = append(message.Replies.Messages, messageReplie...)

	time.Sleep(time.Second * 3)

	return nil
}

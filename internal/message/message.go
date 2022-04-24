package message

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/VladPetriv/tg_scanner/internal/channel"
	"github.com/VladPetriv/tg_scanner/internal/user"
	"github.com/gotd/td/tg"
)

type Message struct {
	ID      int
	Message string
	FromID  user.User
	PeerID  channel.Group
	Replies Replies
	ReplyTo ReplyTo
}

type Replies struct {
	Count    int
	Messages []RepliesMessage
}

type RepliesMessage struct {
	ID      int
	FromID  user.User
	Message string
	ReplyTo interface{}
}

type ReplyTo struct {
	ReplyToMsgID int
}

func GetMessagesFromTelegram(ctx context.Context, data tg.ModifiedMessagesMessages, channelPeer *tg.InputPeerChannel, api *tg.Client) []Message { // nolint
	var msg Message

	result := make([]Message, 0)
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

func GetReplies(ctx context.Context, message *Message, channelPeer *tg.InputPeerChannel, api *tg.Client) (tg.MessagesMessagesClass, error) { // nolint
	bInt := big.NewInt(10000) // nolint

	value, _ := rand.Int(rand.Reader, bInt)

	replies, err := api.MessagesGetReplies(ctx, &tg.MessagesGetRepliesRequest{ // nolint
		Peer:  channelPeer,
		MsgID: message.ID,
		Hash:  value.Int64(),
	})
	if err != nil {
		return nil, fmt.Errorf("error while getting replies:%w", err)
	}

	return replies, nil
}

func ProcessRepliesMessage(ctx context.Context, replies tg.MessagesMessagesClass, cPeer *tg.InputPeerChannel, api *tg.Client) []RepliesMessage {
	repliesMessages := make([]RepliesMessage, 0)

	var replieMessage RepliesMessage

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
			fmt.Printf("error while getting user info for replies[ProcessRepliesMessage]: %s\n", err)

			continue
		}

		replieMessage.FromID = *u
		repliesMessages = append(repliesMessages, replieMessage)
	}

	return repliesMessages
}

func GetIncomingMessages(ctx context.Context, tg_user *tg.User, groups []channel.Group, api *tg.Client) ([]Message, error) {
	msgs := make([]Message, 0)

	var msg Message

	data, err := api.MessagesGetDialogs(ctx, &tg.MessagesGetDialogsRequest{ // nolint
		OffsetPeer: &tg.InputPeerUser{
			UserID:     tg_user.ID,
			AccessHash: tg_user.AccessHash,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error while getting incoming message: %w", err)
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

		// Gettin group info for replie
		for _, group := range groups {
			if msg.PeerID.ChannelID == group.ID {
				msg.PeerID = group
			}
		}

		// Getting user info for replie
		u, err := user.GetUserInfo(ctx, msg.FromID.UserID, msg.ID, &tg.InputPeerChannel{
			ChannelID:  int64(msg.PeerID.ID),
			AccessHash: int64(msg.PeerID.AccessHash),
		}, api)
		if err != nil {
			fmt.Printf("error while getting user info for incoming message: %s", err)

			continue
		}

		msg.FromID = *u

		msgs = append(msgs, msg)
	}

	return msgs, nil
}

func GetRepliesForMessageBeforeSave(ctx context.Context, message *Message, api *tg.Client) error {
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

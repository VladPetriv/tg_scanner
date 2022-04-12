package message

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/VladPetriv/tg_scanner/internal/channel"
	"github.com/gotd/td/tg"
)

type Message struct {
	ID      int
	Message string
	FromID  tg.PeerUser
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
	FromID  tg.PeerUser
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

		repliesMessages := ProcessRepliesMessage(replies)
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

func ProcessRepliesMessage(replies tg.MessagesMessagesClass) []RepliesMessage {
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

		repliesMessages = append(repliesMessages, replieMessage)
	}

	return repliesMessages
}

func GetIncomingMessages(ctx context.Context, user *tg.User, groups []channel.Group, api *tg.Client) ([]Message, error) {
	msgs := make([]Message, 0)

	var msg Message

	data, err := api.MessagesGetDialogs(ctx, &tg.MessagesGetDialogsRequest{ // nolint
		OffsetPeer: &tg.InputPeerUser{
			UserID:     user.ID,
			AccessHash: user.AccessHash,
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
		for _, group := range groups {
			if msg.PeerID.ChannelID == group.ID {
				msg.PeerID = group
			}
		}

		msgs = append(msgs, msg)
	}

	return msgs, nil
}

func GetRepliesForMessageBeforeSave(ctx context.Context, message *Message, api *tg.Client) error {
	replie, err := GetReplies(ctx, message, &tg.InputPeerChannel{
		ChannelID:  int64(message.PeerID.ID),
		AccessHash: int64(message.PeerID.AccessHash),
	}, api)
	if err != nil {
		return err
	}

	messageReplie := ProcessRepliesMessage(replie)

	for _, replie := range messageReplie {
		message.Replies.Messages = append(message.Replies.Messages, replie)
	}

	return nil
}

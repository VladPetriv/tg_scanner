package message

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/gotd/td/tg"
	"github.com/sirupsen/logrus"
)

type Message struct {
	ID      int
	Message string
	FromID  tg.PeerUser
	PeerID  interface{}
	Replies Replies
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

	value, err := rand.Int(rand.Reader, bInt)
	if err != nil {
		logrus.Errorf("ERROR_WHILE_GENERATE_RANDOM_INT:%s", err)
	}

	replies, err := api.MessagesGetReplies(ctx, &tg.MessagesGetRepliesRequest{ // nolint
		Peer:  channelPeer,
		MsgID: message.ID,
		Hash:  value.Int64(),
	})
	if err != nil {
		return nil, fmt.Errorf("ERROR_WHILE_GETTING_REPLIES:%w", err)
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

func GetIncomingMessages(ctx context.Context, user *tg.User, api *tg.Client) ([]Message, error) {
	msgs := make([]Message, 0)

	var msg Message

	data, err := api.MessagesGetDialogs(ctx, &tg.MessagesGetDialogsRequest{ // nolint
		OffsetPeer: &tg.InputPeerUser{
			UserID:     user.ID,
			AccessHash: user.AccessHash,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("ERROR_WHILE_GETTING_DIALOGS:%w", err)
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

		msgs = append(msgs, msg)
	}

	return msgs, nil
}

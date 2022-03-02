package message

import (
	"context"
	"encoding/json"

	"github.com/gotd/td/tg"
)

type Message struct {
	ID      int
	Message string
	FromID  tg.PeerUser
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

func GetMessagesFromTelegram(ctx context.Context, data tg.ModifiedMessagesMessages, api *tg.Client, channelPeer *tg.InputPeerChannel) []Message {
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

		replies, err := GetReplies(ctx, api, channelPeer, &msg)
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

func GetReplies(ctx context.Context, api *tg.Client, channelPeer *tg.InputPeerChannel, message *Message) (tg.MessagesMessagesClass, error) {
	replies, err := api.MessagesGetReplies(ctx, &tg.MessagesGetRepliesRequest{
		Peer:  channelPeer,
		MsgID: message.ID,
		Hash:  324321432513432,
	})
	if err != nil {
		return nil, err
	}

	return replies, nil
}

func ProcessRepliesMessage(replies tg.MessagesMessagesClass) []RepliesMessage {
	var rms []RepliesMessage
	var rm RepliesMessage

	data, _ := replies.AsModified()
	for _, replie := range data.GetMessages() {
		encodedData, err := json.Marshal(replie)
		if err != nil {
			continue
		}
		err = json.Unmarshal(encodedData, &rm)
		if err != nil {
			continue
		}
		rms = append(rms, rm)
	}
	return rms
}

package filter

import (
	"strings"

	"github.com/VladPetriv/tg_scanner/internal/model"
)

func Messages(msg model.TgMessage) bool {
	if msg.ReplyTo.ReplyToMsgID == 0 {
		if strings.Contains(msg.Message, "?") {
			msg.Message = strings.ReplaceAll(msg.Message, "\n", " ")

			return true
		}
	}

	return false
}

func RemoveDuplicatesFromMessages(msgs []model.TgMessage) []model.TgMessage {
	allMessages := make(map[string]bool)
	messages := make([]model.TgMessage, 0)

	for _, m := range msgs {
		if _, status := allMessages[m.Message]; !status {
			allMessages[m.Message] = true

			messages = append(messages, m)
		}
	}

	return messages
}

func RemoveDuplicatesFromReplies(replie *model.TgReplies) {
	if replie.Count == 0 {
		return
	}

	allReplies := make(map[string]bool)
	replies := make([]model.TgRepliesMessage, 0)

	for _, r := range replie.Messages {
		if _, status := allReplies[r.Message]; !status {
			allReplies[r.Message] = true

			replies = append(replies, r)
		}
	}

	replie.Messages = replies
}

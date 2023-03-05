package filter

import (
	"strings"

	"github.com/VladPetriv/tg_scanner/internal/model"
)

func IsQuestion(message string) bool {
	return strings.Contains(message, "?")
}

func ReplaceUnexpectedSymbols(message string) string {
	return strings.ReplaceAll(message, "\n", " ")
}

func RemoveDuplicatesFromMessages(msgs []model.Message) []model.Message {
	allMessages := make(map[string]bool)
	messages := make([]model.Message, 0)

	for _, m := range msgs {
		if _, status := allMessages[m.Message]; !status {
			allMessages[m.Message] = true

			messages = append(messages, m)
		}
	}

	return messages
}

func RemoveDuplicatesFromReplies(reply *model.Replies) {
	if reply.Count == 0 {
		return
	}

	allReplies := make(map[string]bool)
	replies := make([]model.RepliesMessage, 0)

	for _, r := range reply.Messages {
		if _, status := allReplies[r.Message]; !status {
			allReplies[r.Message] = true

			replies = append(replies, r)
		}
	}

	reply.Messages = replies
}

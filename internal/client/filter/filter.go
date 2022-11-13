package filter

import (
	"strings"

	"github.com/VladPetriv/tg_scanner/internal/model"
)

// ProcessMessage check if message is question and remove unexpected symbols from message
func ProcessMessage(msg *model.TgMessage) bool {
	isQuestion := checkIfMessageIsQuestion(msg)

	if isQuestion {
		msg.Message = replaceUnexpectedSymbols(msg.Message)
	}

	return isQuestion
}

func checkIfMessageIsQuestion(msg *model.TgMessage) bool {
	return strings.Contains(msg.Message, "?")
}

func replaceUnexpectedSymbols(message string) string {
	return strings.ReplaceAll(message, "\n", " ")
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

func RemoveDuplicatesFromReplies(reply *model.TgReplies) {
	if reply.Count == 0 {
		return
	}

	allReplies := make(map[string]bool)
	replies := make([]model.TgRepliesMessage, 0)

	for _, r := range reply.Messages {
		if _, status := allReplies[r.Message]; !status {
			allReplies[r.Message] = true

			replies = append(replies, r)
		}
	}

	reply.Messages = replies
}

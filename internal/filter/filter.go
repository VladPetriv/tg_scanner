package filter

import (
	"strings"

	"github.com/VladPetriv/tg_scanner/internal/model"
)

func Messages(msg *model.TgMessage) (*model.TgMessage, bool) {
	if msg.ReplyTo.ReplyToMsgID == 0 {
		if strings.Contains(msg.Message, "?") {
			msg.Message = strings.ReplaceAll(msg.Message, "\n", " ")

			return msg, true
		}
	}

	return nil, false
}

func RemoveDuplicateByMessage(msgs []model.TgMessage) []model.TgMessage {
	allKeys := make(map[string]bool)
	messages := make([]model.TgMessage, 0)

	for _, item := range msgs {
		if _, value := allKeys[item.Message]; !value {
			allKeys[item.Message] = true

			messages = append(messages, item)
		}
	}

	return messages
}

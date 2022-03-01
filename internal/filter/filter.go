package filter

import (
	"strings"

	"example.com/test/m/internal/message"
)

func FilterMessages(msg *message.Message) (*message.Message, bool) {
	if strings.Contains(msg.Message, "?") {
		msg.Message = strings.ReplaceAll(msg.Message, "\n", " ")
		return msg, true
	}
	return nil, false
}

func RemoveDuplicateByMessage(msgs []message.Message) []message.Message {
	allKeys := make(map[string]bool)
	messages := []message.Message{}
	for _, item := range msgs {
		if _, value := allKeys[item.Message]; !value {
			allKeys[item.Message] = true
			messages = append(messages, item)
		}
	}
	return messages
}

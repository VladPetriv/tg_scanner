package model

import (
	"fmt"

	"github.com/cnf/structhash"
)

type TgMessage struct {
	ID         int
	Message    string
	FromID     TgUser    `json:"FromID"`
	PeerID     TgGroup   `json:"PeerID"`
	Replies    TgReplies `json:"Replies"`
	ReplyTo    TgReplyTo `json:"ReplyTo"`
	Media      Media     `json:"Media"`
	MessageURL string
	ImageURL   string
}

func (m TgMessage) GetHash() (string, error) {
	hash, err := structhash.Hash(m, 1)
	if err != nil {
		return "", fmt.Errorf("hash message data: %w", err)
	}

	return hash, nil
}

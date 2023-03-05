package model

import (
	"fmt"

	"github.com/cnf/structhash"
)

type Message struct {
	ID         int     `json:"ID"`
	Message    string  `json:"Message"`
	FromID     User    `json:"FromID"`
	PeerID     Group   `json:"PeerID"`
	Replies    Replies `json:"Replies"`
	ReplyTo    ReplyTo `json:"ReplyTo"`
	Media      Media   `json:"Media"`
	MessageURL string  `json:"MessageURL"`
	ImageURL   string  `json:"ImageURL"`
}

func (m Message) GetHash() (string, error) {
	hash, err := structhash.Hash(m, 1)
	if err != nil {
		return "", fmt.Errorf("hash message data: %w", err)
	}

	return hash, nil
}

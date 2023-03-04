package model

import (
	"fmt"

	"github.com/cnf/structhash"
	"github.com/gotd/td/tg"
)

type TgGroup struct {
	ID         int64
	ChannelID  int64
	Title      string
	AccessHash int64
	Username   string
	Photo      tg.ChatPhoto
	ImageURL   string
}

func (g TgGroup) GetHash() (string, error) {
	hash, err := structhash.Hash(g, 1)
	if err != nil {
		return "", fmt.Errorf("hash group data: %w", err)
	}

	return hash, nil
}

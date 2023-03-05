package model

import (
	"fmt"

	"github.com/cnf/structhash"
)

type Group struct {
	ID         int64      `json:"ID"`
	ChannelID  int64      `json:"ChannelID"`
	Title      string     `json:"Title"`
	AccessHash int64      `json:"AccessHash"`
	Username   string     `json:"Username"`
	Photo      GroupPhoto `json:"Photo"`
	ImageURL   string     `json:"ImageURL"`
}

type GroupPhoto struct {
	PhotoID int64 `json:"PhotoID"`
}

func (g Group) GetHash() (string, error) {
	hash, err := structhash.Hash(g, 1)
	if err != nil {
		return "", fmt.Errorf("hash group data: %w", err)
	}

	return hash, nil
}

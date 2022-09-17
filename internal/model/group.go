package model

import "github.com/gotd/td/tg"

type TgGroup struct {
	ID         int64
	ChannelID  int64
	Title      string
	AccessHash int64
	Username   string
	Photo      tg.ChatPhoto
	ImageURL   string
}

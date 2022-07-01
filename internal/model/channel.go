package model

import "github.com/gotd/td/tg"

type Channel struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Title    string `json:"title"`
	ImageURL string `json:"imageUrl"`
}

type TgChannel struct {
	ID         int64
	ChannelID  int64
	Title      string
	AccessHash int64
	Username   string
	Photo      tg.ChatPhoto
	Image      *Image
}

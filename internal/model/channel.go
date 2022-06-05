package model

import "github.com/gotd/td/tg"

type Channel struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Title    string `json:"title"`
	PhotoURL string `json:"photoUrl"`
}

type TgChannel struct {
	ID         int
	ChannelID  int
	Title      string
	AccessHash int
	Username   string
	Photo      tg.ChatPhoto
	Image      *Image
}

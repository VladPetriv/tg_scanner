package model

import "github.com/gotd/td/tg"

type Channel struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	Title    string `db:"title"`
	ImageURL string `db:"imageurl"`
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

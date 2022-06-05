package model

import (
	"github.com/gotd/td/tg"
)

type Message struct {
	ID           int    `json:"id"`
	ChannelID    int    `json:"channelId"`
	UserID       int    `json:"userId"`
	Title        string `json:"title"`
	MessageURL   string `json:"messageURL"`
	RepliesCount int
}

type TgMessage struct {
	ID      int
	Message string
	FromID  TgUser
	PeerID  TgChannel
	Replies TgReplies
	ReplyTo TgReplyTo
	Media   tg.MessageMediaPhoto
	Image   *Image
}

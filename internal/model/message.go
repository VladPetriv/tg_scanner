package model

type Message struct {
	ID        int    `json:"id"`
	ChannelID int    `json:"channelId`
	Title     string `json:"title"`
}

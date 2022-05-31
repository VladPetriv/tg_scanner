package model

type Message struct {
	ID           int    `json:"id"`
	ChannelID    int    `json:"channelId"`
	UserID       int    `json:"userId"`
	Title        string `json:"title"`
	MessageURL   string `json:"messageURL"`
	RepliesCount int
}

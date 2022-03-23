package model

type Message struct {
	Id        int    `json:"Id"`
	ChannelId int    `json:"ChannelId"`
	Title     string `json:"Title"`
}

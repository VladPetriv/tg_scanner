package model

type Message struct {
	ID           int    `json:"id"`
	ChannelID    int    `json:"channelId"`
	UserID       int    `json:"userId"`
	Title        string `json:"title"`
	MessageURL   string `json:"messageURL"`
	Image        string `json:"image"`
	RepliesCount int
}

type TgMessage struct {
	ID      int
	Message string
	FromID  TgUser    `json:"FromID"`
	PeerID  TgChannel `json:"PeerID"`
	Replies TgReplies `json:"Replies"`
	ReplyTo TgReplyTo `json:"ReplyTo"`
	Media   Media     `json:"Media"`
	Image   *Image
}

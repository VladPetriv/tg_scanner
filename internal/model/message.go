package model

type Message struct {
	ID           int    `db:"id"`
	ChannelID    int    `db:"channel_id"`
	UserID       int    `db:"user_id"`
	Title        string `db:"title"`
	MessageURL   string `db:"message_url"`
	ImageURL     string `db:"imageurl"`
	RepliesCount int    `db:"count"`
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

package model

type TgMessage struct {
	ID         int
	Message    string
	FromID     TgUser    `json:"FromID"`
	PeerID     TgChannel `json:"PeerID"`
	Replies    TgReplies `json:"Replies"`
	ReplyTo    TgReplyTo `json:"ReplyTo"`
	Media      Media     `json:"Media"`
	MessageURL string
	ImageURL   string
}

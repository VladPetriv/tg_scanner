package model

type Message struct {
	ID         int
	Message    string
	FromID     User    `json:"FromID"`
	PeerID     Group   `json:"PeerID"`
	Replies    Replies `json:"Replies"`
	ReplyTo    ReplyTo `json:"ReplyTo"`
	Media      Media   `json:"Media"`
	MessageURL string
	ImageURL   string
}

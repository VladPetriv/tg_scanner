package model

type Replies struct {
	Count    int
	Messages []RepliesMessage
}

type RepliesMessage struct {
	ID       int
	FromID   TgUser
	Message  string
	ReplyTo  interface{}
	Media    Media
	ImageURL string
}

type ReplyTo struct {
	ReplyToMsgID int
}

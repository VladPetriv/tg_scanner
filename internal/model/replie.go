package model

type TgReplies struct {
	Count    int
	Messages []TgRepliesMessage
}

type TgRepliesMessage struct {
	ID       int
	FromID   TgUser
	Message  string
	ReplyTo  interface{}
	Media    Media
	ImageURL string
}

type TgReplyTo struct {
	ReplyToMsgID int
}

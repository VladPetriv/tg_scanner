package model

type Replie struct {
	ID        int    `json:"id"`
	MessageID int    `json:"messageId"`
	UserID    int    `json:"userId"`
	Title     string `json:"title"`
}

type TgReplies struct {
	Count    int
	Messages []TgRepliesMessage
}

type TgRepliesMessage struct {
	ID      int
	FromID  TgUser
	Message string
	ReplyTo interface{}
}

type TgReplyTo struct {
	ReplyToMsgID int
}

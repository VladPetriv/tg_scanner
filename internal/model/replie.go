package model

type Replie struct {
	ID        int    `db:"id"`
	MessageID int    `db:"message_id"`
	UserID    int    `db:"user_id"`
	Title     string `db:"title"`
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

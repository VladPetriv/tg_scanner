package model

type Replie struct {
	ID        int    `db:"id"`
	MessageID int    `db:"message_id"`
	UserID    int    `db:"user_id"`
	Title     string `db:"title"`
	ImageURL  string `db:"imageurl"`
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
	Media   Media
}

type TgReplyTo struct {
	ReplyToMsgID int
}

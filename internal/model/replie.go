package model

type Replies struct {
	Count    int              `json:"Count"`
	Messages []RepliesMessage `json:"Messages"`
}

type RepliesMessage struct {
	ID       int    `json:"ID"`
	FromID   User   `json:"FromID"`
	Message  string `json:"Message"`
	Media    Media  `json:"Media"`
	ImageURL string `json:"ImageURL"`
}

type ReplyTo struct {
	ReplyToMsgID int `json:"ReplyToMsgID"`
}

package model

type Replie struct {
	ID        int    `json:"id"`
	MessageID int    `json:"messageId"`
	UserID    int    `json:"userId"`
	Title     string `json:"title"`
}

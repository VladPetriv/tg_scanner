package model

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Username  string `json:"username"`
	MessageID string `json:"messageId"`
	PhotoURL  string `json:"photoUrl"`
}

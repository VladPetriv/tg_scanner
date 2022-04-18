package model

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	FullName string `json:"fullname"`
	PhotoURL string `json:"photoUrl"`
}

package model

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	FullName string `json:"fullName"`
	PhotoURL string `json:"photoUrl"`
}

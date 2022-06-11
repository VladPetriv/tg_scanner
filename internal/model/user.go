package model

import "github.com/gotd/td/tg"

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	FullName string `json:"fullName"`
	PhotoURL string `json:"photoUrl"`
}

type TgUser struct {
	ID         int
	UserID     int
	FirstName  string
	LastName   string
	Username   string
	Photo      tg.UserProfilePhoto
	AccessHash int
	Image      *Image
}

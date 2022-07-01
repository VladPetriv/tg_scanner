package model

import "github.com/gotd/td/tg"

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	FullName string `json:"fullName"`
	ImageURL string `json:"imageUrl"`
}

type TgUser struct {
	ID         int64
	UserID     int64
	AccessHash int64
	FirstName  string
	LastName   string
	Username   string
	Photo      tg.UserProfilePhoto
	Image      *Image
}

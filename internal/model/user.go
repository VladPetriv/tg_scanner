package model

import "github.com/gotd/td/tg"

type TgUser struct {
	ID         int64
	UserID     int64
	AccessHash int64
	FirstName  string
	LastName   string
	Username   string
	Photo      tg.UserProfilePhoto
	ImageURL   string
}

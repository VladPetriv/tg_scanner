package model

type User struct {
	ID         int64
	UserID     int64
	AccessHash int64
	FirstName  string
	LastName   string
	Username   string
	Photo      UserProfilePhoto
	ImageURL   string
	Fullname   string
}

type UserProfilePhoto struct {
	PhotoID int64
}

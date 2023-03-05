package model

type User struct {
	ID         int64     `json:"ID"`
	UserID     int64     `json:"UserID"`
	AccessHash int64     `json:"AccessHash"`
	FirstName  string    `json:"FirstName"`
	LastName   string    `json:"LastName"`
	Username   string    `json:"Username"`
	Photo      UserPhoto `json:"Photo"`
	ImageURL   string    `json:"ImageURL"`
	Fullname   string    `json:"FullName"`
}

type UserPhoto struct {
	PhotoID int64 `json:"PhotoID"`
}

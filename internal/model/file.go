package model

import "github.com/gotd/td/tg"

type Image struct {
	Bytes []byte
}
type Media struct {
	Photo *Photo `json:"Photo"`
}

type Photo struct {
	ID            int
	AccessHash    int
	FileReference []byte
	Sizes         []tg.PhotoSize
}

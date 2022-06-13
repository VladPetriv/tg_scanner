package model

import "github.com/gotd/td/tg"

type Image struct {
	Bytes []byte
}
type Media struct {
	Photo *Photo `json:"Photo"`
}

type Photo struct {
	ID            int64
	AccessHash    int64
	FileReference []byte
	Sizes         []tg.PhotoSize
}

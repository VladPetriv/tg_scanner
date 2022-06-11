package model

import "github.com/gotd/td/tg"

type Image struct {
	Bytes []byte
}

type Media struct {
	Photo *tg.Photo `json:"Photo"`
}

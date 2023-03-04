package model

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
	Sizes         []PhotoSize
}

type PhotoSize struct {
	Type string
	W    int
	H    int
	Size int
}

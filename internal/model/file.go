package model

type Image struct {
	Bytes []byte
}
type Media struct {
	Photo *Photo `json:"Photo"`
}

type Photo struct {
	ID            int64       `json:"ID"`
	AccessHash    int64       `json:"AccessHash"`
	FileReference []byte      `json:"FileReference"`
	Sizes         []PhotoSize `json:"Sizes"`
}

type PhotoSize struct {
	Type string `json:"Type"`
	W    int    `json:"W"`
	H    int    `json:"H"`
	Size int    `json:"Size"`
}

func (p PhotoSize) GetType() string {
	return p.Type
}

package photo

import (
	"context"

	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/gotd/td/tg"
)

type Photo interface {
	DecodePhoto(photo tg.UploadFileClass) (*model.Image, error)
	CreatePhoto(img *model.Image, name string) (string, error)
	ProcessPhoto(ctx context.Context, photoData tg.UploadFileClass, name string) (string, error)
}

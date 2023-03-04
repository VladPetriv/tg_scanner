package photo

import (
	"context"

	"github.com/gotd/td/tg"

	"github.com/VladPetriv/tg_scanner/internal/model"
)

type Photo interface {
	ProcessPhoto(ctx context.Context, photoData tg.UploadFileClass, name string) (string, error)
	decodePhoto(photo tg.UploadFileClass) (*model.Image, error)
	createPhoto(img *model.Image, name string) (string, error)
}

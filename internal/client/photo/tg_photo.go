package photo

import (
	"context"

	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/gotd/td/tg"
)

type TgPhoto interface {
	ProcessPhoto(ctx context.Context, photoData tg.UploadFileClass, name string) (string, error)
	decodePhoto(photo tg.UploadFileClass) (*model.Image, error)
	createPhoto(img *model.Image, name string) (string, error)
}

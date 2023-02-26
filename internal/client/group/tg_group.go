package group

import (
	"context"

	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/gotd/td/tg"
)

type TgGroup interface {
	GetGroups(ctx context.Context) ([]model.TgGroup, error)
	GetGroupPhoto(ctx context.Context, group model.TgGroup) (tg.UploadFileClass, error)
	CreateFilesForGroups(groups []model.TgGroup) error
}

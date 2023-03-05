package group

import (
	"context"

	"github.com/gotd/td/tg"

	"github.com/VladPetriv/tg_scanner/internal/model"
)

type Group interface {
	GetGroups(ctx context.Context) ([]model.Group, error)
	GetGroupPhoto(ctx context.Context, group model.Group) (tg.UploadFileClass, error)
	CreateFilesForGroups(groups []model.Group) error
}

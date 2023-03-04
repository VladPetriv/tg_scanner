package group

import (
	"context"

	"github.com/gotd/td/tg"

	"github.com/VladPetriv/tg_scanner/internal/model"
)

type Group interface {
	GetGroups(ctx context.Context) ([]model.TgGroup, error)
	GetMessagesFromGroupHistory(ctx context.Context, groupPeer *tg.InputPeerChannel) (tg.MessagesMessagesClass, error)
	GetGroupPhoto(ctx context.Context, group model.TgGroup) (tg.UploadFileClass, error)
	CreateFilesForGroups(groups []model.TgGroup) error
}

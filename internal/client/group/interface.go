package group

import (
	"context"

	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/gotd/td/tg"
)

type Group interface {
	GetGroups(ctx context.Context) ([]model.TgGroup, error)
	GetMessagesFromGroupHistory(ctx context.Context, groupPeer *tg.InputPeerChannel) (tg.MessagesMessagesClass, error)
	GetGroupPhoto(ctx context.Context, group *model.TgGroup) (tg.UploadFileClass, error)
}

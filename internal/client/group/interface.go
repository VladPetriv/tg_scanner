package group

import (
	"context"

	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/gotd/td/tg"
)

type Group interface {
	GetMessagesFromGroupHistory(ctx context.Context, groupPeer *tg.InputPeerChannel) (tg.MessagesMessagesClass, error)
	GetGroups(ctx context.Context) ([]model.TgGroup, error)
	GetGroupPhoto(ctx context.Context, group *model.TgGroup) (tg.UploadFileClass, error)
}

package user

import (
	"context"

	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/gotd/td/tg"
)

type TgUser interface {
	GetUser(ctx context.Context, data interface{}, groupPeer *tg.InputPeerChannel) (*model.TgUser, error)
	GetUserPhoto(ctx context.Context, user model.TgUser) (tg.UploadFileClass, error)
}

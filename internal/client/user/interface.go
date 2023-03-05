package user

import (
	"context"

	"github.com/gotd/td/tg"

	"github.com/VladPetriv/tg_scanner/internal/model"
)

type User interface {
	GetUser(ctx context.Context, data interface{}, groupPeer *tg.InputPeerChannel) (*model.User, error)
	GetUserPhoto(ctx context.Context, user model.User) (tg.UploadFileClass, error)
}

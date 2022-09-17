package message

import (
	"context"

	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/gotd/td/tg"
)

type Message interface {
	ProcessHistoryMessages(ctx context.Context, data tg.ModifiedMessagesMessages, groupPeer *tg.InputPeerChannel) []model.TgMessage
	ProcessIncomingMessages(ctx context.Context, tgUser *tg.User, groups []model.TgGroup) ([]model.TgMessage, error)
	GetMessagePhoto(ctx context.Context, message model.TgMessage) (tg.UploadFileClass, error)
	CheckMessagePhotoStatus(ctx context.Context, message *model.TgMessage) (bool, error)
}

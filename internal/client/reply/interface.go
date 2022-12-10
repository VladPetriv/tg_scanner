package reply

import (
	"context"

	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/gotd/td/tg"
)

type Reply interface {
	GetReplies(ctx context.Context, message model.TgMessage, groupPeer *tg.InputPeerChannel) (tg.MessagesMessagesClass, error)
	ParseTelegramReplies(ctx context.Context, replies tg.MessagesMessagesClass, groupPeer *tg.InputPeerChannel) []model.TgRepliesMessage
	GetReplyPhoto(ctx context.Context, reply model.TgRepliesMessage) (tg.UploadFileClass, error)
}

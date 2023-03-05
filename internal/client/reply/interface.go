package reply

import (
	"context"

	"github.com/gotd/td/tg"

	"github.com/VladPetriv/tg_scanner/internal/model"
)

type Reply interface {
	GetReplies(ctx context.Context, msg model.Message, groupPeer *tg.InputPeerChannel) (tg.MessagesMessagesClass, error)
	ParseTelegramReplies(ctx context.Context, replies tg.MessagesMessagesClass, groupPeer *tg.InputPeerChannel) []model.RepliesMessage //nolint:lll
	GetReplyPhoto(ctx context.Context, reply model.RepliesMessage) (tg.UploadFileClass, error)
}

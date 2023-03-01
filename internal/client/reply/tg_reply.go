package reply

import (
	"context"

	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/gotd/td/tg"
)

type TgReply interface {
	GetReplies(ctx context.Context, msg model.TgMessage) (tg.MessagesMessagesClass, error)
	ParseTelegramReplies(ctx context.Context, replies tg.MessagesMessagesClass, groupPeer *tg.InputPeerChannel) []model.TgRepliesMessage //nolint:lll
	GetReplyPhoto(ctx context.Context, reply model.TgRepliesMessage) (tg.UploadFileClass, error)
}

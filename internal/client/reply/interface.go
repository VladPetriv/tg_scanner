package reply

import (
	"context"

	"github.com/gotd/td/tg"

	"github.com/VladPetriv/tg_scanner/internal/model"
)

type Reply interface {
	GetReplies(ctx context.Context, message model.Message) ([]model.RepliesMessage, error)
	GetReplyPhoto(ctx context.Context, reply model.RepliesMessage) (tg.UploadFileClass, error)
}

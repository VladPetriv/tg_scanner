package message

import (
	"context"

	"github.com/gotd/td/tg"

	"github.com/VladPetriv/tg_scanner/internal/model"
)

type Message interface {
	ParseHistoryMessages(ctx context.Context, data tg.ModifiedMessagesMessages, groupPeer *tg.InputPeerChannel) []model.TgMessage //nolint:lll
	ParseIncomingMessages(ctx context.Context, tgUser tg.User, groups []model.TgGroup) ([]model.TgMessage, error)
	GetMessagePhoto(ctx context.Context, message model.TgMessage) (tg.UploadFileClass, error)
	CheckMessagePhotoStatus(ctx context.Context, message *model.TgMessage) (bool, error)
	WriteMessagesToFile(messages []model.TgMessage, fileName string)
	GetMessagesFromFile(pathToFile string) ([]model.TgMessage, error)
}

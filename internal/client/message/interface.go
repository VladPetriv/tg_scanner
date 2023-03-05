package message

import (
	"context"

	"github.com/gotd/td/tg"

	"github.com/VladPetriv/tg_scanner/internal/model"
)

type Message interface {
	ParseHistoryMessages(ctx context.Context, data tg.ModifiedMessagesMessages, groupPeer *tg.InputPeerChannel) []model.Message //nolint:lll
	ParseIncomingMessages(ctx context.Context, tgUser tg.User, groups []model.Group) ([]model.Message, error)
	GetMessagePhoto(ctx context.Context, message model.Message) (tg.UploadFileClass, error)
	CheckMessagePhotoStatus(ctx context.Context, message *model.Message) (bool, error)
	WriteMessagesToFile(messages []model.Message, fileName string)
	GetMessagesFromFile(pathToFile string) ([]model.Message, error)
}

package message

import (
	"context"

	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/gotd/td/tg"
)

type Message interface {
	// GetQuestionsFromGroupHistory retrieves all questions asked in a group's history.
	GetQuestionsFromGroupHistory(ctx context.Context, groupPeer *tg.InputPeerChannel) ([]model.TgMessage, error)
	ParseIncomingMessages(ctx context.Context, tgUser tg.User, groups []model.TgGroup) ([]model.TgMessage, error)
	GetMessagePhoto(ctx context.Context, message model.TgMessage) (tg.UploadFileClass, error)
	CheckMessagePhotoStatus(ctx context.Context, message *model.TgMessage) (bool, error)
	WriteMessagesToFile(messages []model.TgMessage, fileName string)
	GetMessagesFromFile(pathToFile string) ([]model.TgMessage, error)
}

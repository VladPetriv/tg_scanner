package message

import (
	"context"

	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/gotd/td/tg"
)

type TgMessage interface {
	// GetQuestionsFromGroupHistory retrieves all questions asked in a group's history.
	GetQuestionsFromGroupHistory(ctx context.Context, groupPeer *tg.InputPeerChannel) ([]model.TgMessage, error)
	// GetIncomingQuestionsFromGroup retrieves incoming questions from the given Telegram groups
	// that are directed to the specified user.
	GetIncomingQuestionsFromGroup(ctx context.Context, tgUser tg.User) ([]model.TgMessage, error)
	GetMessagePhoto(ctx context.Context, message model.TgMessage) (tg.UploadFileClass, error)
	CheckMessagePhotoStatus(ctx context.Context, message *model.TgMessage) (bool, error)
	WriteMessagesToFile(messages []model.TgMessage, fileName string)
	GetMessagesFromFile(pathToFile string) ([]model.TgMessage, error)
}

package client

import (
	"context"

	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/gotd/td/tg"
)

type AppClient interface {
	ProcessMessagesFromGroupHistory(groups []model.TgGroup)
	GetIncomingMessages(user tg.User, groups []model.TgGroup)
	ValidateAndPushGroupsToQueue(ctx context.Context) ([]model.TgGroup, error)
	PushMessagesToQueue()
}

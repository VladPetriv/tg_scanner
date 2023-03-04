package client

import (
	"context"

	"github.com/gotd/td/tg"

	"github.com/VladPetriv/tg_scanner/internal/model"
)

type AppClient interface {
	GetHistoryMessages(groups []model.TgGroup)
	GetIncomingMessages(user tg.User, groups []model.TgGroup)
	ValidateAndPushGroupsToQueue(ctx context.Context) ([]model.TgGroup, error)
	PushMessagesToQueue()
}

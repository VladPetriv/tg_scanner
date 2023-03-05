package client

import (
	"context"

	"github.com/gotd/td/tg"

	"github.com/VladPetriv/tg_scanner/internal/model"
)

type AppClient interface {
	GetHistoryMessages(groups []model.Group)
	GetIncomingMessages(user tg.User, groups []model.Group)
	ValidateAndPushGroupsToQueue(ctx context.Context) ([]model.Group, error)
	PushMessagesToQueue()
}

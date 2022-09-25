package client

import (
	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/gotd/td/tg"
)

type AppClient interface {
	GetHistoryMessages(groups []model.TgGroup)
	GetIncomingMessages(user *tg.User, groups []model.TgGroup)
	PushToQueue()
}

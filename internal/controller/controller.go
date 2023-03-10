package controller

import "github.com/VladPetriv/tg_scanner/internal/model"

type Controller interface {
	SendGroupData(group model.Group) error
	SendMessageData(message model.Message) error
}

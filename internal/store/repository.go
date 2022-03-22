package store

import (
	"github.com/VladPetriv/tg_scanner/internal/model"
)

type ChannelRepo interface {
	GetChannels() (*[]model.Channel, error)
	GetChannel(channelId int) (*model.Channel, error)
	GetChannelByName(name string) (*model.Channel, error)
	CreateChannel(channel *model.Channel) error
	DeleteChannel(channelId int) error
}

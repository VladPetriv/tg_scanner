package store

import (
	"github.com/VladPetriv/tg_scanner/internal/model"
)

type ChannelRepo interface {
	GetChannel(channelID int) (*model.Channel, error)
	GetChannelByName(name string) (*model.Channel, error)
	CreateChannel(channel *model.Channel) (int, error)
}

type MessageRepo interface {
	GetMessage(messageID int) (*model.Message, error)
	GetMessageByName(name string) (*model.Message, error)
	CreateMessage(message *model.Message) (int, error)
}

type ReplieRepo interface {
	GetReplie(replieID int) (*model.Replie, error)
	GetReplieByName(name string) (*model.Replie, error)
	CreateReplie(replie *model.Replie) (int, error)
}

type UserRepo interface {
	GetUsers() ([]model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	CreateUser(user *model.User) (int, error)
}

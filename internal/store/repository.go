package store

import (
	"github.com/VladPetriv/tg_scanner/internal/model"
)

type ChannelRepo interface {
	GetChannels() ([]model.Channel, error)
	GetChannel(channelID int) (*model.Channel, error)
	GetChannelByName(name string) (*model.Channel, error)
	CreateChannel(channel *model.Channel) (int, error)
	DeleteChannel(channelID int) error
}

type MessageRepo interface {
	GetMessages() ([]model.Message, error)
	GetMessage(messageID int) (*model.Message, error)
	GetMessageByName(name string) (*model.Message, error)
	CreateMessage(message *model.Message) (int, error)
	DeleteMessage(messageID int) error
}

type ReplieRepo interface {
	GetReplies() ([]model.Replie, error)
	GetReplie(replieID int) (*model.Replie, error)
	GetReplieByName(name string) (*model.Replie, error)
	CreateReplie(replie *model.Replie) (int, error)
	DeleteReplie(replieID int) error
}

type UserRepo interface {
	GetUsers() ([]model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	CreateUser(user *model.User) (int, error)
}

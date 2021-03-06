package store

import (
	"github.com/VladPetriv/tg_scanner/internal/model"
)

//go:generate mockery --dir . --name ChannelRepo --output ./mocks
type ChannelRepo interface {
	GetChannelByName(name string) (*model.Channel, error)
	CreateChannel(channel *model.Channel) (int, error)
}

//go:generate mockery --dir . --name MessageRepo --output ./mocks
type MessageRepo interface {
	GetMessageByName(name string) (*model.Message, error)
	CreateMessage(message *model.Message) (int, error)
	DeleteMessageByID(messageID int) (int, error)
	GetMessagesWithRepliesCount() ([]model.Message, error)
}

//go:generate mockery --dir . --name ReplieRepo --output ./mocks
type ReplieRepo interface {
	GetReplieByName(name string) (*model.Replie, error)
	CreateReplie(replie *model.Replie) (int, error)
}

//go:generate mockery --dir . --name UserRepo --output ./mocks
type UserRepo interface {
	GetUserByUsername(username string) (*model.User, error)
	CreateUser(user *model.User) (int, error)
}

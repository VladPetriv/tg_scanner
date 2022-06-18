package service

import "github.com/VladPetriv/tg_scanner/internal/model"

//go:generate mockery --dir . --name ChannelService --output ./mocks
type ChannelService interface {
	GetChannelByName(name string) (*model.Channel, error)
	CreateChannel(channel *model.Channel) error
}

//go:generate mockery --dir . --name MessageService --output ./mocks
type MessageService interface {
	GetMessageByName(name string) (*model.Message, error)
	CreateMessage(message *model.Message) (int, error)
	DeleteMessageByID(messageID int) error
	GetMessagesWithRepliesCount() ([]model.Message, error)
}

//go:generate mockery --dir . --name ReplieService --output ./mocks
type ReplieService interface {
	GetReplieByName(name string) (*model.Replie, error)
	CreateReplie(replie *model.Replie) error
}

//go:generate mockery --dir . --name UserService --output ./mocks
type UserService interface {
	GetUserByUsername(username string) (*model.User, error)
	CreateUser(user *model.User) (int, error)
}

package service

import "github.com/VladPetriv/tg_scanner/internal/model"

//go:generate mockery --dir . --name ChannelService --output ./mocks
type ChannelService interface {
	GetChannel(channelID int) (*model.Channel, error)
	GetChannelByName(name string) (*model.Channel, error)
	CreateChannel(channel *model.Channel) error
}

//go:generate mockery --dir . --name MessageService --output ./mocks
type MessageService interface {
	GetMessage(messagelID int) (*model.Message, error)
	GetMessageByName(name string) (*model.Message, error)
	CreateMessage(message *model.Message) (int, error)
}

//go:generate mockery --dir . --name ReplieService --output ./mocks
type ReplieService interface {
	GetReplie(replieID int) (*model.Replie, error)
	GetReplieByName(name string) (*model.Replie, error)
	CreateReplie(replie *model.Replie) error
}

//go:generate mockery --dir . --name UserService --output ./mocks
type UserService interface {
	GetUserByUsername(username string) (*model.User, error)
	CreateUser(user *model.User) (int, error)
}

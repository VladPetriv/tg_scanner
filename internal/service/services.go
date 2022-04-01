package service

import "github.com/VladPetriv/tg_scanner/internal/model"

type ChannelService interface {
	GetChannels() ([]model.Channel, error)
	GetChannel(channelId int) (*model.Channel, error)
	GetChannelByName(name string) (*model.Channel, error)
	CreateChannel(channel *model.Channel) error
	DeleteChannel(channelId int) error
}
type MessageService interface {
	GetMessages() (*[]model.Message, error)
	GetMessage(messagelId int) (*model.Message, error)
	GetMessageByName(name string) (*model.Message, error)
	CreateMessage(message *model.Message) error
	DeleteMessage(messageId int) error
}

type ReplieService interface {
	GetReplies() (*[]model.Replie, error)
	GetReplie(replieId int) (*model.Replie, error)
	GetReplieByName(name string) (*model.Replie, error)
	CreateReplie(replie *model.Replie) error
	DeleteReplie(replieId int) error
}

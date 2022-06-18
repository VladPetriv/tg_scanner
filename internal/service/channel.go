package service

import (
	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/internal/store"
	"github.com/VladPetriv/tg_scanner/pkg/utils"
)

type ChannelDBService struct {
	store *store.Store
}

func NewChannelDBService(store *store.Store) *ChannelDBService {
	return &ChannelDBService{
		store: store,
	}
}

func (s *ChannelDBService) GetChannelByName(name string) (*model.Channel, error) {
	channel, err := s.store.Channel.GetChannelByName(name)
	if err != nil {
		return nil, &utils.ServiceError{
			ServiceName:       "Channel",
			ServiceMethodName: "GetChannelByName",
			ErrorValue:        err,
		}
	}

	if channel == nil {
		return nil, nil
	}

	return channel, nil
}

func (s *ChannelDBService) CreateChannel(channel *model.Channel) error {
	candidate, err := s.GetChannelByName(channel.Name)
	if err != nil {
		return err
	}

	if candidate != nil {
		return &utils.RecordIsExistError{RecordName: "channel", Name: channel.Name}
	}

	_, err = s.store.Channel.CreateChannel(channel)
	if err != nil {
		return &utils.ServiceError{
			ServiceName:       "Channel",
			ServiceMethodName: "CreateChannel",
			ErrorValue:        err,
		}
	}

	return nil
}

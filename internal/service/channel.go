package service

import (
	"fmt"

	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/internal/store"
)

type ChannelDbService struct {
	store *store.Store
}

func NewChannelDbService(store *store.Store) *ChannelDbService {
	return &ChannelDbService{
		store: store,
	}
}

func (s *ChannelDbService) GetChannels() ([]model.Channel, error) {
	channels, err := s.store.Channel.GetChannels()
	if err != nil {
		return nil, fmt.Errorf("[Channel] Service.GetChannels error: %w", err)
	}

	if channels == nil {
		return nil, fmt.Errorf("channels not found")
	}

	return channels, nil
}

func (s *ChannelDbService) GetChannel(channelId int) (*model.Channel, error) {
	channel, err := s.store.Channel.GetChannel(channelId)
	if err != nil {
		return nil, fmt.Errorf("[Channel] Service.GetChannel error: %w", err)
	}

	if channel == nil {
		return nil, fmt.Errorf("channel not found")
	}

	return channel, nil
}

func (s *ChannelDbService) GetChannelByName(name string) (*model.Channel, error) {
	channel, err := s.store.Channel.GetChannelByName(name)
	if err != nil {
		return nil, fmt.Errorf("[Channel] Service.GetChannelByName error: %w", err)
	}
	if channel == nil {
		return nil, fmt.Errorf("channel not found")
	}

	return channel, nil
}

func (s *ChannelDbService) CreateChannel(channel *model.Channel) error {
	candidate, err := s.GetChannelByName(channel.Name)
	if err != nil {
		return err
	}

	if candidate.Name == channel.Name {
		return fmt.Errorf("channel with name %s is exist", channel.Name)
	}

	_, err = s.store.Channel.CreateChannel(channel)
	if err != nil {
		return fmt.Errorf("[Channel] Service error: %w", err)
	}

	return nil
}

func (s *ChannelDbService) DeleteChannel(channelId int) error {
	err := s.store.Channel.DeleteChannel(channelId)
	if err != nil {
		return fmt.Errorf("[Channel] Service error: %w", err)
	}

	return nil
}

package service

import (
	"fmt"

	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/internal/store"
)

type ChannelDBService struct {
	store *store.Store
}

func NewChannelDBService(store *store.Store) *ChannelDBService {
	return &ChannelDBService{
		store: store,
	}
}

func (s *ChannelDBService) GetChannels() ([]model.Channel, error) {
	channels, err := s.store.Channel.GetChannels()
	if err != nil {
		return nil, fmt.Errorf("[Channel] Service.GetChannels error: %w", err)
	}

	if channels == nil {
		return nil, fmt.Errorf("channels not found")
	}

	return channels, nil
}

func (s *ChannelDBService) GetChannel(channelID int) (*model.Channel, error) {
	channel, err := s.store.Channel.GetChannel(channelID)
	if err != nil {
		return nil, fmt.Errorf("[Channel] Service.GetChannel error: %w", err)
	}

	if channel == nil {
		return nil, fmt.Errorf("channel not found")
	}

	return channel, nil
}

func (s *ChannelDBService) GetChannelByName(name string) (*model.Channel, error) {
	channel, err := s.store.Channel.GetChannelByName(name)
	if err != nil {
		return nil, fmt.Errorf("[Channel] Service.GetChannelByName error: %w", err)
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
		return fmt.Errorf("channel with name %s is exist", channel.Name)
	}

	_, err = s.store.Channel.CreateChannel(channel)
	if err != nil {
		return fmt.Errorf("[Channel] Service error: %w", err)
	}

	return nil
}

func (s *ChannelDBService) DeleteChannel(channelID int) error {
	err := s.store.Channel.DeleteChannel(channelID)
	if err != nil {
		return fmt.Errorf("[Channel] Service error: %w", err)
	}

	return nil
}

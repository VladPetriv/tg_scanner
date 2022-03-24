package service

import (
	"fmt"

	"github.com/VladPetriv/tg_scanner/internal/store"
)

type Manager struct {
	Channel ChannelService
	Message MessageService
}

func NewManager(store *store.Store) (*Manager, error) {
	if store == nil {
		return nil, fmt.Errorf("No store provided")
	}

	return &Manager{
		Channel: NewChannelDbService(store),
		Message: NewMessageDbService(store),
	}, nil
}
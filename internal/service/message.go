package service

import (
	"fmt"

	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/internal/store"
)

type MessageDbService struct {
	store *store.Store
}

func NewMessageDbService(store *store.Store) *MessageDbService {
	return &MessageDbService{
		store: store,
	}
}

func (s *MessageDbService) GetMessages() ([]model.Message, error) {
	messages, err := s.store.Message.GetMessages()
	if err != nil {
		return nil, fmt.Errorf("[Message] Service.GetMessages error: %w", err)
	}

	if messages == nil {
		return nil, nil
	}

	return messages, nil
}

func (s *MessageDbService) GetMessage(messageId int) (*model.Message, error) {
	message, err := s.store.Message.GetMessage(messageId)
	if err != nil {
		return nil, fmt.Errorf("[Message] Service.GetMessage error: %w", err)
	}

	if message == nil {
		return nil, fmt.Errorf("message not found")
	}

	return message, nil
}

func (s *MessageDbService) GetMessageByName(name string) (*model.Message, error) {
	message, err := s.store.Message.GetMessageByName(name)
	if err != nil {
		return nil, fmt.Errorf("[Message] Service.GetMessageByName error: %w", err)
	}

	if message == nil {
		return nil, nil
	}

	return message, nil
}

func (s *MessageDbService) CreateMessage(message *model.Message) error {
	candidate, err := s.GetMessageByName(message.Title)
	if err != nil {
		return err
	}

	if candidate != nil {
		return fmt.Errorf("message with name %s is exist", message.Title)
	}

	_, err = s.store.Message.CreateMessage(message)
	if err != nil {
		return fmt.Errorf("[Message] Service.CreateMessage error: %w", err)
	}

	return nil
}

func (s *MessageDbService) DeleteMessage(messageId int) error {
	err := s.store.Message.DeleteMessage(messageId)
	if err != nil {
		return fmt.Errorf("[Message] Service.DeleteMessage error: %w", err)
	}

	return nil
}

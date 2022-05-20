package service

import (
	"fmt"

	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/internal/store"
)

type MessageDBService struct {
	store *store.Store
}

func NewMessageDBService(store *store.Store) *MessageDBService {
	return &MessageDBService{
		store: store,
	}
}

func (s *MessageDBService) GetMessage(messageID int) (*model.Message, error) {
	message, err := s.store.Message.GetMessage(messageID)
	if err != nil {
		return nil, fmt.Errorf("[Message] Service.GetMessage error: %w", err)
	}

	return message, nil
}

func (s *MessageDBService) GetMessageByName(name string) (*model.Message, error) {
	message, err := s.store.Message.GetMessageByName(name)
	if err != nil {
		return nil, fmt.Errorf("[Message] Service.GetMessageByName error: %w", err)
	}

	if message == nil {
		return nil, nil
	}

	return message, nil
}

func (s *MessageDBService) CreateMessage(message *model.Message) (int, error) {
	candidate, err := s.GetMessageByName(message.Title)
	if err != nil {
		return 0, err
	}

	if candidate != nil {
		return candidate.ID, fmt.Errorf("message with name %s is exist", message.Title)
	}

	id, err := s.store.Message.CreateMessage(message)
	if err != nil {
		return id, fmt.Errorf("[Message] Service.CreateMessage error: %w", err)
	}

	return id, nil
}

func (s *MessageDBService) DeleteMessageByID(messageID int) error {
	_, err := s.store.Message.DeleteMessageByID(messageID)
	if err != nil {
		return fmt.Errorf("[Message] Service.DeleteMessageByID error: %w", err)
	}

	return nil
}

func (s *MessageDBService) GetMessagesWithRepliesCount() ([]model.Message, error) {
	messages, err := s.store.Message.GetMessagesWithRepliesCount()
	if err != nil {
		return nil, fmt.Errorf("[Message] Service.GetMessagesWithRepliesCount error: %w", err)
	}

	if messages == nil {
		return nil, fmt.Errorf("messages not found")
	}

	return messages, nil
}

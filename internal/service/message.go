package service

import (
	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/internal/store"
	"github.com/VladPetriv/tg_scanner/pkg/utils"
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
		return nil, &utils.ServiceError{
			ServiceName:       "Message",
			ServiceMethodName: "GetMessage",
			ErrorValue:        err,
		}
	}

	if message == nil {
		return nil, &utils.NotFoundError{Name: "message"}
	}

	return message, nil
}

func (s *MessageDBService) GetMessageByName(name string) (*model.Message, error) {
	message, err := s.store.Message.GetMessageByName(name)
	if err != nil {
		return nil, &utils.ServiceError{
			ServiceName:       "Message",
			ServiceMethodName: "GetMessageByName",
			ErrorValue:        err,
		}
	}

	if message == nil {
		return nil, nil
	}

	return message, nil
}

func (s *MessageDBService) CreateMessage(message *model.Message) (int, error) {
	candidate, err := s.store.Message.GetMessageByName(message.Title)
	if err != nil {
		return 0, &utils.ServiceError{
			ServiceName:       "Message",
			ServiceMethodName: "CreateMessage [GetMessageByName]",
			ErrorValue:        err,
		}
	}

	if candidate != nil && candidate.ChannelID == message.ChannelID {
		return candidate.ID, &utils.RecordIsExistError{RecordName: "message", Name: candidate.Title}
	}

	id, err := s.store.Message.CreateMessage(message)
	if err != nil {
		return id, &utils.ServiceError{
			ServiceName:       "Message",
			ServiceMethodName: "CreateMessage",
			ErrorValue:        err,
		}
	}

	return id, nil
}

func (s *MessageDBService) DeleteMessageByID(messageID int) error {
	_, err := s.store.Message.DeleteMessageByID(messageID)
	if err != nil {
		return &utils.ServiceError{
			ServiceName:       "Message",
			ServiceMethodName: "DeleteMessageByID",
			ErrorValue:        err,
		}
	}

	return nil
}

func (s *MessageDBService) GetMessagesWithRepliesCount() ([]model.Message, error) {
	messages, err := s.store.Message.GetMessagesWithRepliesCount()
	if err != nil {
		return nil, &utils.ServiceError{
			ServiceName:       "Message",
			ServiceMethodName: "GetMessagesWithRepliesCount",
			ErrorValue:        err,
		}
	}

	if messages == nil {
		return nil, &utils.NotFoundError{Name: "messages"}
	}

	return messages, nil
}

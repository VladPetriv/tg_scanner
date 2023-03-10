package controller

import (
	"fmt"

	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/pkg/logger"
	"github.com/VladPetriv/tg_scanner/pkg/queue"
)

type messageBroker struct {
	queue  queue.Queue
	logger *logger.Logger
}

var _ Controller = (*messageBroker)(nil)

func New(queue queue.Queue, logger *logger.Logger) Controller {
	return &messageBroker{
		queue:  queue,
		logger: logger,
	}
}

func (m messageBroker) SendGroupData(data model.Group) error {
	logger := m.logger

	err := m.queue.SendMessage("group/create", data)
	if err != nil {
		logger.Error().Err(err).Msg("send group data to queue")
		return fmt.Errorf("send group data to queue: %w", err)
	}

	logger.Info().Msg("successfully send group data to queue")
	return nil
}

func (m messageBroker) SendMessageData(data model.Message) error {
	logger := m.logger

	err := m.queue.SendMessage("message/create", data)
	if err != nil {
		logger.Error().Err(err).Msg("send message data to queue")
		return fmt.Errorf("send message data to queue: %w", err)
	}

	logger.Info().Msg("successfully sent message data to queue")
	return nil
}

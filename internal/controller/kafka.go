package controller

import (
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/VladPetriv/tg_scanner/pkg/config"
)

type queue struct {
	cfg *config.Config
}

func New(cfg *config.Config) *queue {
	return &queue{
		cfg: cfg,
	}
}

var _ Controller = (*queue)(nil)

func (q queue) PushDataToQueue(topic string, data interface{}) error {
	producer, err := connectAsProducer(q.cfg.KafkaAddr)
	if err != nil {
		return fmt.Errorf("connect as producer error: %w", err)
	}

	defer producer.Close()

	encodedData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal data error: %w", err)
	}

	queueMessage := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(encodedData),
	}

	_, _, err = producer.SendMessage(queueMessage)
	if err != nil {
		return fmt.Errorf("send message into kafka error: %w", err)
	}

	return nil
}

func connectAsProducer(addr string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Retry.Max = 5

	conn, err := sarama.NewSyncProducer([]string{addr}, config)
	if err != nil {
		return nil, fmt.Errorf("create producer error: %w", err)
	}

	return conn, nil
}

package queue

import (
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"
)

type kafka struct {
	Address string
}

var _ Queue = (*kafka)(nil)

func New(address string) Queue {
	return &kafka{
		Address: address,
	}
}

func (k kafka) SendMessage(topic string, data interface{}) error {
	producer, err := connectAsProducer(k.Address)
	if err != nil {
		return fmt.Errorf("connect as producer to queue: %w", err)
	}

	defer producer.Close()

	encodedData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal data: %w", err)
	}

	queueMessage := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(encodedData),
	}

	_, _, err = producer.SendMessage(queueMessage)
	if err != nil {
		return fmt.Errorf("send message to queue: %w", err)
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
		return nil, fmt.Errorf("create new producer: %w", err)
	}

	return conn, nil
}

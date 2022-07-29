package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/VladPetriv/tg_scanner/internal/model"
)

func connectAsProducer(addr string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Retry.Max = 5

	conn, err := sarama.NewSyncProducer([]string{addr}, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	return conn, nil
}

func PushDataToQueue[T model.TgMessage | model.TgChannel](topic, addr string, data T) error {
	producer, err := connectAsProducer(addr)
	if err != nil {
		return fmt.Errorf("failed to connect as producer: %w", err)
	}

	defer producer.Close()

	encodedData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	queueMessage := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(encodedData),
	}

	_, _, err = producer.SendMessage(queueMessage)
	if err != nil {
		return fmt.Errorf("failed send message into kafka: %w", err)
	}

	return nil
}

package queue

type Queue interface {
	SendMessageToQueue(topic string, message interface{}) error
}

package queue

type Queue interface {
	SendMessage(topic string, data interface{}) error
}

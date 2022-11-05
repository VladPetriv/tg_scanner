package controller

type Controller interface {
	PushDataToQueue(topic string, data interface{}) error
}

package utils

import (
	"fmt"

	"github.com/VladPetriv/tg_scanner/logger"
)

func CheckError(err error, log *logger.Logger) {
	if err == nil {
		return
	}

	switch t := err.(type) {
	case *RecordIsExistError:
		log.Warn(t)
	case *NotFoundError, *ServiceError, *CreateError, *GettingError:
		log.Error(t)
	default:
		log.Info(t)
	}
}

type NotFoundError struct {
	Name string
}

func (n *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", n.Name)
}

type RecordIsExistError struct {
	RecordName string
	Name       string
}

func (r *RecordIsExistError) Error() string {
	return fmt.Sprintf("%s with name %s is exist", r.RecordName, r.Name)
}

type ServiceError struct {
	ServiceName       string
	ServiceMethodName string
	ErrorValue        error
}

func (s *ServiceError) Error() string {
	return fmt.Sprintf("[%s] Service.%s error: %s", s.ServiceName, s.ServiceMethodName, s.ErrorValue)
}

type CreateError struct {
	Name       string
	ErrorValue error
}

func (c *CreateError) Error() string {
	return fmt.Sprintf("create %s error: %s", c.Name, c.ErrorValue)
}

type GettingError struct {
	Name       string
	ErrorValue error
}

func (g *GettingError) Error() string {
	return fmt.Sprintf("get %s error: %s", g.Name, g.ErrorValue)
}

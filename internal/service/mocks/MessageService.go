// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"
import model "github.com/VladPetriv/tg_scanner/internal/model"

// MessageService is an autogenerated mock type for the MessageService type
type MessageService struct {
	mock.Mock
}

// CreateMessage provides a mock function with given fields: message
func (_m *MessageService) CreateMessage(message *model.Message) (int, error) {
	ret := _m.Called(message)

	var r0 int
	if rf, ok := ret.Get(0).(func(*model.Message) int); ok {
		r0 = rf(message)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.Message) error); ok {
		r1 = rf(message)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMessage provides a mock function with given fields: messagelID
func (_m *MessageService) GetMessage(messagelID int) (*model.Message, error) {
	ret := _m.Called(messagelID)

	var r0 *model.Message
	if rf, ok := ret.Get(0).(func(int) *model.Message); ok {
		r0 = rf(messagelID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Message)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(messagelID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMessageByName provides a mock function with given fields: name
func (_m *MessageService) GetMessageByName(name string) (*model.Message, error) {
	ret := _m.Called(name)

	var r0 *model.Message
	if rf, ok := ret.Get(0).(func(string) *model.Message); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Message)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

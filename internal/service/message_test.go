package service_test

import (
	"errors"
	"testing"

	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/internal/service"
	"github.com/VladPetriv/tg_scanner/internal/store"
	"github.com/VladPetriv/tg_scanner/internal/store/mocks"
	"github.com/stretchr/testify/assert"
)

func TestMessageService_GetMessage(t *testing.T) {
	tests := []struct {
		name    string
		mock    func(messageRepo *mocks.MessageRepo)
		input   int
		want    *model.Message
		wantErr bool
		err     error
	}{
		{
			name: "Ok: [Message found]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessage", 1).Return(&model.Message{ID: 1, UserID: 1, ChannelID: 1, Title: "test"}, nil)
			},
			input: 1,
			want:  &model.Message{ID: 1, UserID: 1, ChannelID: 1, Title: "test"},
		},
		{
			name: "Error: [Message not found]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessage", 404).Return(nil, errors.New("message not found"))
			},
			input:   404,
			wantErr: true,
			err:     errors.New("[Message] Service.GetMessage error: message not found"),
		},
		{
			name: "Error: [Store error]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessage", 1).Return(nil, errors.New("error while getting message: some error"))
			},
			input:   1,
			wantErr: true,
			err:     errors.New("[Message] Service.GetMessage error: error while getting message: some error"),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		messageRepo := &mocks.MessageRepo{}
		messageService := service.NewMessageDBService(&store.Store{Message: messageRepo})
		tt.mock(messageRepo)

		got, err := messageService.GetMessage(tt.input)
		if tt.wantErr {
			assert.Error(t, err)
			assert.Equal(t, tt.err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		messageRepo.AssertExpectations(t)
	}
}

func TestMessageService_CreateMessage(t *testing.T) {
	input := &model.Message{
		ChannelID: 1,
		UserID:    1,
		Title:     "test",
	}

	tests := []struct {
		name    string
		mock    func(messageRepo *mocks.MessageRepo)
		input   *model.Message
		want    int
		wantErr bool
		err     error
	}{
		{
			name: "Ok: [Message created]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessageByName", input.Title).Return(nil, nil)
				messageRepo.On("CreateMessage", input).Return(1, nil)
			},
			input: input,
			want:  1,
		},
		{
			name: "Error: [Message is exist]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessageByName", input.Title).Return(input, nil)
			},
			input:   input,
			want:    0,
			wantErr: true,
			err:     errors.New("message with name test is exist"),
		},
		{
			name: "Error: [Store error]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessageByName", input.Title).Return(nil, errors.New("error while getting message: some error"))
			},
			input:   input,
			want:    0,
			wantErr: true,
			err:     errors.New("[Message] Service.GetMessageByName error: error while getting message: some error"),
		},
		{
			name: "Error: [Store error]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessageByName", input.Title).Return(nil, nil)
				messageRepo.On("CreateMessage", input).Return(0, errors.New("error while creating message: some error"))
			},
			input:   input,
			want:    0,
			wantErr: true,
			err:     errors.New("[Message] Service.CreateMessage error: error while creating message: some error"),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		messageRepo := *&mocks.MessageRepo{}
		messageService := service.NewMessageDBService(&store.Store{Message: &messageRepo})
		tt.mock(&messageRepo)

		got, err := messageService.CreateMessage(tt.input)
		if tt.wantErr {
			assert.Error(t, err)
			assert.Equal(t, tt.err.Error(), err.Error())
			assert.Equal(t, tt.want, got)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		messageRepo.AssertExpectations(t)
	}
}

func TestMessageService_GetMessageByName(t *testing.T) {
	tests := []struct {
		name    string
		mock    func(messageRepo *mocks.MessageRepo)
		input   string
		want    *model.Message
		wantErr bool
		err     error
	}{
		{
			name: "Ok: [Message found]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessageByName", "test").Return(&model.Message{ID: 1, ChannelID: 1, UserID: 1, Title: "test"}, nil)

			},
			input: "test",
			want:  &model.Message{ID: 1, ChannelID: 1, UserID: 1, Title: "test"},
		},
		{
			name: "Error: [Message not found]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessageByName", "test").Return(nil, nil)
			},
			input: "test",
			want:  nil,
		},
		{
			name: "Error: [Store error]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessageByName", "test").Return(nil, errors.New("error while getting message: some error"))
			},
			input:   "test",
			wantErr: true,
			err:     errors.New("[Message] Service.GetMessageByName error: error while getting message: some error"),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		messageRepo := &mocks.MessageRepo{}
		messageService := service.NewMessageDBService(&store.Store{Message: messageRepo})
		tt.mock(messageRepo)

		got, err := messageService.GetMessageByName(tt.input)
		if tt.wantErr {
			assert.Error(t, err)
			assert.Equal(t, tt.err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		messageRepo.AssertExpectations(t)
	}
}

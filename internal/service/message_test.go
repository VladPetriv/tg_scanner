package service_test

import (
	"fmt"
	"testing"

	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/internal/service"
	"github.com/VladPetriv/tg_scanner/internal/store"
	"github.com/VladPetriv/tg_scanner/internal/store/mocks"
	"github.com/VladPetriv/tg_scanner/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestMessageService_CreateMessage(t *testing.T) {
	input := &model.Message{
		ChannelID: 1,
		UserID:    1,
		Title:     "test",
		ImageURL:  "test.jpg",
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
				messageRepo.On("GetMessageByName", input.Title).
					Return(nil, nil)

				messageRepo.On("CreateMessage", input).
					Return(1, nil)
			},
			input: input,
			want:  1,
		},
		{
			name: "Error: [Message is exist]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessageByName", input.Title).
					Return(input, nil)
			},
			input:   input,
			want:    0,
			wantErr: true,
			err:     &utils.RecordIsExistError{RecordName: "message", Name: "test"},
		},
		{
			name: "Error: [Store error: 'get message by name']",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessageByName", input.Title).
					Return(nil, &utils.GettingError{Name: "message by name", ErrorValue: fmt.Errorf("some error")})
			},
			input:   input,
			want:    0,
			wantErr: true,
			err: &utils.ServiceError{
				ServiceName:       "Message",
				ServiceMethodName: "GetMessageByName",
				ErrorValue:        fmt.Errorf("get message by name error: some error"),
			},
		},
		{
			name: "Error: [Store error: 'create message']",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessageByName", input.Title).
					Return(nil, nil)

				messageRepo.On("CreateMessage", input).
					Return(0, &utils.CreateError{Name: "message", ErrorValue: fmt.Errorf("some error")})
			},
			input:   input,
			want:    0,
			wantErr: true,
			err: &utils.ServiceError{
				ServiceName:       "Message",
				ServiceMethodName: "CreateMessage",
				ErrorValue:        fmt.Errorf("create message error: some error"),
			},
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		messageRepo := &mocks.MessageRepo{}
		messageService := service.NewMessageDBService(&store.Store{Message: messageRepo})
		tt.mock(messageRepo)

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
	data := &model.Message{
		ID:        1,
		ChannelID: 1,
		UserID:    1,
		Title:     "test",
		ImageURL:  "test.jpg",
	}

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
				messageRepo.On("GetMessageByName", "test").
					Return(data, nil)
			},
			input: "test",
			want:  &model.Message{ID: 1, ChannelID: 1, UserID: 1, Title: "test", ImageURL: "test.jpg"},
		},
		{
			name: "Error: [Message not found]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessageByName", "test").
					Return(nil, nil)
			},
			input: "test",
			want:  nil,
		},
		{
			name: "Error: [Store error]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessageByName", "test").
					Return(nil, &utils.GettingError{Name: "message by name", ErrorValue: fmt.Errorf("some error")})
			},
			input:   "test",
			wantErr: true,
			err: &utils.ServiceError{
				ServiceName:       "Message",
				ServiceMethodName: "GetMessageByName",
				ErrorValue:        fmt.Errorf("get message by name error: some error"),
			},
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

func TestMessageService_GetMessagesWithReplies(t *testing.T) {
	messages := []model.Message{
		{ID: 1, RepliesCount: 2},
		{ID: 2, RepliesCount: 13},
	}

	tests := []struct {
		name    string
		mock    func(messageRepo *mocks.MessageRepo)
		want    []model.Message
		wantErr bool
		err     error
	}{
		{
			name: "Ok: [Messages found]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessagesWithRepliesCount").
					Return(messages, nil)
			},
			want: messages,
		},
		{
			name: "Error: [Messages not found]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessagesWithRepliesCount").
					Return(nil, nil)
			},
			wantErr: true,
			err:     &utils.NotFoundError{Name: "messages with replies count"},
		},
		{
			name: "Error: [Store error]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessagesWithRepliesCount").
					Return(nil, &utils.GettingError{Name: "messages with replies count", ErrorValue: fmt.Errorf("some error")})
			},
			wantErr: true,
			err: &utils.ServiceError{
				ServiceName:       "Message",
				ServiceMethodName: "GetMessagesWithRepliesCount",
				ErrorValue:        fmt.Errorf("get messages with replies count error: some error"),
			},
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		messageRepo := &mocks.MessageRepo{}
		messageService := service.NewMessageDBService(&store.Store{Message: messageRepo})
		tt.mock(messageRepo)

		got, err := messageService.GetMessagesWithRepliesCount()
		if tt.wantErr {
			assert.Error(t, err)
			assert.Equal(t, tt.err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}
	}
}

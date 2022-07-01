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

func TestChannelService_CreateChannel(t *testing.T) {
	input := &model.Channel{
		Name:     "test",
		Title:    "test",
		ImageURL: "test.jpg",
	}

	tests := []struct {
		name    string
		mock    func(channelRepo *mocks.ChannelRepo)
		input   *model.Channel
		wantErr bool
		err     error
	}{
		{
			name: "Ok: [Channel created]",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelByName", input.Name).
					Return(nil, nil)

				channelRepo.On("CreateChannel", input).
					Return(1, nil)
			},
			input: input,
		},

		{
			name: "Error: [Channel is exist]",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelByName", input.Name).
					Return(input, nil)
			},
			input:   input,
			wantErr: true,
			err:     &utils.RecordIsExistError{RecordName: "channel", Name: "test"},
		},
		{
			name: "Error: [Store error: 'get channel by name']",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelByName", input.Name).
					Return(nil, &utils.GettingError{Name: "channel by name", ErrorValue: fmt.Errorf("some error")})
			},
			input:   input,
			wantErr: true,
			err: &utils.ServiceError{
				ServiceName:       "Channel",
				ServiceMethodName: "GetChannelByName",
				ErrorValue:        fmt.Errorf("get channel by name error: some error"),
			},
		},
		{
			name: "Error: [Store error: 'create channel']",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelByName", input.Name).
					Return(nil, nil)

				channelRepo.On("CreateChannel", input).
					Return(0, &utils.CreateError{Name: "channel", ErrorValue: fmt.Errorf("some error")})
			},
			input:   input,
			wantErr: true,
			err: &utils.ServiceError{
				ServiceName:       "Channel",
				ServiceMethodName: "CreateChannel",
				ErrorValue:        fmt.Errorf("create channel error: some error"),
			},
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		channelRepo := &mocks.ChannelRepo{}
		channelService := service.NewChannelDBService(&store.Store{Channel: channelRepo})
		tt.mock(channelRepo)

		err := channelService.CreateChannel(tt.input)
		if tt.wantErr {
			assert.Error(t, err)
			assert.Equal(t, tt.err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, nil, err)
		}

		channelRepo.AssertExpectations(t)
	}
}

func TestChannelService_GetChannelByName(t *testing.T) {
	tests := []struct {
		name    string
		mock    func(channelRepo *mocks.ChannelRepo)
		input   string
		want    *model.Channel
		wantErr bool
		err     error
	}{
		{
			name: "Ok: [Channel found]",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelByName", "test").
					Return(&model.Channel{ID: 1, Name: "test", Title: "test", ImageURL: "test.jpg"}, nil)
			},
			input: "test",
			want:  &model.Channel{ID: 1, Name: "test", Title: "test", ImageURL: "test.jpg"},
		},
		{
			name: "Error: [Channel not found]",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelByName", "test").
					Return(nil, nil)
			},
			input: "test",
			want:  nil,
		},
		{
			name: "Error [Store error]",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelByName", "test").
					Return(nil, &utils.GettingError{Name: "channel by name", ErrorValue: fmt.Errorf("some error")})
			},
			input:   "test",
			wantErr: true,
			err: &utils.ServiceError{
				ServiceName:       "Channel",
				ServiceMethodName: "GetChannelByName",
				ErrorValue:        fmt.Errorf("get channel by name error: some error"),
			},
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		channelRepo := mocks.ChannelRepo{}
		channelService := service.NewChannelDBService(&store.Store{Channel: &channelRepo})
		tt.mock(&channelRepo)

		got, err := channelService.GetChannelByName(tt.input)
		if tt.wantErr {
			assert.Error(t, err)
			assert.Equal(t, tt.err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		channelRepo.AssertExpectations(t)
	}
}

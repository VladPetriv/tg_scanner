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

func TestChannelService_GetChannel(t *testing.T) {
	tests := []struct {
		name         string
		expectations func(channelRepo *mocks.ChannelRepo)
		input        int
		want         *model.Channel
		wantErr      bool
		err          error
	}{
		{
			name: "Ok: [Channel found]",
			expectations: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannel", 1).Return(
					&model.Channel{ID: 1, Name: "test", Title: "test", PhotoURL: "test.jpg"}, nil,
				)
			},
			input: 1,
			want:  &model.Channel{ID: 1, Name: "test", Title: "test", PhotoURL: "test.jpg"},
		},
		{
			name: "Error: [channel not found]",
			expectations: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannel", 404).Return(nil, nil)
			},
			input:   404,
			err:     errors.New("channel not found"),
			wantErr: true,
		},
		{
			name: "Error: [Store error]",
			expectations: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannel", 0).Return(nil, errors.New("error while getting channel: some error"))
			},
			input:   0,
			err:     errors.New("[Channel] Service.GetChannel error: error while getting channel: some error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		channelRepo := &mocks.ChannelRepo{}
		channelService := service.NewChannelDBService(&store.Store{Channel: channelRepo})
		tt.expectations(channelRepo)

		got, err := channelService.GetChannel(tt.input)
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

func TestChannelService_CreateChannel(t *testing.T) {
	input := &model.Channel{
		Name:     "test",
		Title:    "test",
		PhotoURL: "test.jpg",
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
				channelRepo.On("GetChannelByName", input.Name).Return(nil, nil)
				channelRepo.On("CreateChannel", input).Return(1, nil)
			},
			input: input,
		},

		{
			name: "Error: [Channel is exist]",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelByName", input.Name).Return(input, nil)
			},
			input:   input,
			wantErr: true,
			err:     errors.New("channel with name test is exist"),
		},
		{
			name: "Error: [Store error]",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelByName", input.Name).Return(nil, errors.New("error while getting channel: some error"))
			},
			input:   input,
			wantErr: true,
			err:     errors.New("[Channel] Service.GetChannelByName error: error while getting channel: some error"),
		},
		{
			name: "Error: [Store error]",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelByName", input.Name).Return(nil, nil)
				channelRepo.On("CreateChannel", input).Return(0, errors.New("error while creating channel: some error"))
			},
			input:   input,
			wantErr: true,
			err:     errors.New("[Channel] Service.CreateChannel error: error while creating channel: some error"),
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
				channelRepo.On("GetChannelByName", "test").Return(&model.Channel{ID: 1, Name: "test", Title: "test", PhotoURL: "test.jpg"}, nil)
			},
			input: "test",
			want:  &model.Channel{ID: 1, Name: "test", Title: "test", PhotoURL: "test.jpg"},
		},
		{
			name: "Error: [Channel not found]",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelByName", "test").Return(nil, nil)
			},
			input: "test",
			want:  nil,
		},
		{
			name: "Error [Store error]",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelByName", "test").Return(nil, errors.New("error while getting channel: some error"))
			},
			input:   "test",
			wantErr: true,
			err:     errors.New("[Channel] Service.GetChannelByName error: error while getting channel: some error"),
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

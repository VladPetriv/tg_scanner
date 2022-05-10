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

func TestReplieService_GetReplie(t *testing.T) {
	tests := []struct {
		name    string
		mock    func(replieRepo *mocks.ReplieRepo)
		input   int
		want    *model.Replie
		wantErr bool
		err     error
	}{
		{
			name: "Ok: [Replie found]",
			mock: func(replieRepo *mocks.ReplieRepo) {
				replieRepo.On("GetReplie", 1).Return(&model.Replie{ID: 1, MessageID: 1, UserID: 1, Title: "test"}, nil)
			},
			input: 1,
			want:  &model.Replie{ID: 1, MessageID: 1, UserID: 1, Title: "test"},
		},
		{
			name: "Error: [Replie not found]",
			mock: func(replieRepo *mocks.ReplieRepo) {
				replieRepo.On("GetReplie", 404).Return(nil, errors.New("replie not found"))
			},
			input:   404,
			wantErr: true,
			err:     errors.New("[Replie] Service.GetReplie error: replie not found"),
		},
		{
			name: "Error: [Store error]",
			mock: func(replieRepo *mocks.ReplieRepo) {
				replieRepo.On("GetReplie", 1).Return(nil, errors.New("error while getting replie: some error"))
			},
			input:   1,
			wantErr: true,
			err:     errors.New("[Replie] Service.GetReplie error: error while getting replie: some error"),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		replieRepo := &mocks.ReplieRepo{}
		replieService := service.NewReplieDBService(&store.Store{Replie: replieRepo})
		tt.mock(replieRepo)

		got, err := replieService.GetReplie(tt.input)
		if tt.wantErr {
			assert.Error(t, err)
			assert.Equal(t, tt.err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		replieRepo.AssertExpectations(t)
	}
}

func TestReplieService_CreateReplie(t *testing.T) {
	input := &model.Replie{
		UserID:    1,
		MessageID: 1,
		Title:     "test",
	}
	tests := []struct {
		name    string
		mock    func(replieRepo *mocks.ReplieRepo)
		input   *model.Replie
		want    int
		wantErr bool
		err     error
	}{
		{
			name: "Ok: [Replie created]",
			mock: func(replieRepo *mocks.ReplieRepo) {
				replieRepo.On("GetReplieByName", input.Title).Return(nil, nil)
				replieRepo.On("CreateReplie", input).Return(1, nil)
			},
			input: input,
		},
		{
			name: "Error: [Replie is exist]",
			mock: func(replieRepo *mocks.ReplieRepo) {
				replieRepo.On("GetReplieByName", input.Title).Return(input, nil)
			},
			input:   input,
			wantErr: true,
			err:     errors.New("replie with name test is exist"),
		},
		{
			name: "Error: [Store error]",
			mock: func(replieRepo *mocks.ReplieRepo) {
				replieRepo.On("GetReplieByName", input.Title).Return(nil, errors.New("error while getting replie: some error"))
			},
			input:   input,
			wantErr: true,
			err:     errors.New("[Replie] Service.GetReplieByName error: error while getting replie: some error"),
		},
		{
			name: "Error: [Store error]",
			mock: func(replieRepo *mocks.ReplieRepo) {
				replieRepo.On("GetReplieByName", input.Title).Return(nil, nil)
				replieRepo.On("CreateReplie", input).Return(0, errors.New("error while creating replie: some error"))
			},
			input:   input,
			wantErr: true,
			err:     errors.New("[Replie] Service.CreateReplie error: error while creating replie: some error"),
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		replieRepo := &mocks.ReplieRepo{}
		replieService := service.NewReplieDBService(&store.Store{Replie: replieRepo})
		tt.mock(replieRepo)

		err := replieService.CreateReplie(tt.input)
		if tt.wantErr {
			assert.Error(t, err)
			assert.Equal(t, tt.err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
		}

		replieRepo.AssertExpectations(t)
	}
}

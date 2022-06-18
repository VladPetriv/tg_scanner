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
				replieRepo.On("GetReplieByName", input.Title).
					Return(nil, nil)

				replieRepo.On("CreateReplie", input).
					Return(1, nil)
			},
			input: input,
		},
		{
			name: "Error: [Store error]",
			mock: func(replieRepo *mocks.ReplieRepo) {
				replieRepo.On("GetReplieByName", input.Title).
					Return(nil, &utils.GettingError{Name: "replie by name", ErrorValue: fmt.Errorf("some error")})
			},
			input:   input,
			wantErr: true,
			err: &utils.ServiceError{
				ServiceName:       "Replie",
				ServiceMethodName: "GetReplieByName",
				ErrorValue:        fmt.Errorf("get replie by name error: some error"),
			},
		},
		{
			name: "Error: [Store error]",
			mock: func(replieRepo *mocks.ReplieRepo) {
				replieRepo.On("GetReplieByName", input.Title).
					Return(nil, nil)

				replieRepo.On("CreateReplie", input).
					Return(0, &utils.CreateError{Name: "replie", ErrorValue: fmt.Errorf("some error")})
			},
			input:   input,
			wantErr: true,
			err: &utils.ServiceError{
				ServiceName:       "Replie",
				ServiceMethodName: "CreateReplie",
				ErrorValue:        fmt.Errorf("create replie error: some error"),
			},
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

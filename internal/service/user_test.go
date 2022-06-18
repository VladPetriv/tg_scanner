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

func TestUserService_GetUserByUsername(t *testing.T) {
	tests := []struct {
		name    string
		mock    func(userRepo *mocks.UserRepo)
		input   string
		want    *model.User
		wantErr bool
		err     error
	}{
		{
			name: "Ok: [User found]",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUserByUsername", "test").
					Return(&model.User{ID: 1, Username: "test", FullName: "test test", PhotoURL: "test.jpg"}, nil)
			},
			input: "test",
			want:  &model.User{ID: 1, Username: "test", FullName: "test test", PhotoURL: "test.jpg"},
		},
		{
			name: "Error [User not found]",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUserByUsername", "test").
					Return(nil, nil)
			},
			input: "test",
			want:  nil,
		},
		{
			name: "Error: [Store error]",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUserByUsername", "test").
					Return(nil, &utils.GettingError{Name: "user", ErrorValue: fmt.Errorf("some error")})
			},
			input:   "test",
			wantErr: true,
			err: &utils.ServiceError{
				ServiceName:       "User",
				ServiceMethodName: "GetUserByUsername",
				ErrorValue:        fmt.Errorf("get user error: some error"),
			},
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		userRepo := &mocks.UserRepo{}
		userService := service.NewUserDBService(&store.Store{User: userRepo})
		tt.mock(userRepo)

		got, err := userService.GetUserByUsername(tt.input)
		if tt.wantErr {
			assert.Error(t, err)
			assert.Equal(t, tt.err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		userRepo.AssertExpectations(t)
	}
}

func TestUserService_CreateUser(t *testing.T) {
	input := &model.User{
		Username: "test",
		FullName: "test test",
		PhotoURL: "test.jpg",
	}

	tests := []struct {
		name    string
		mock    func(userRepo *mocks.UserRepo)
		input   *model.User
		want    int
		wantErr bool
		err     error
	}{
		{
			name: "Ok: [User created]",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUserByUsername", input.Username).
					Return(nil, nil)

				userRepo.On("CreateUser", input).
					Return(1, nil)
			},
			input: input,
			want:  1,
		},
		{
			name: "Error: [User is exist]",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUserByUsername", input.Username).
					Return(input, nil)
			},
			input:   input,
			wantErr: true,
			err:     &utils.RecordIsExistError{RecordName: "user", Name: "test"},
		},
	}

	for _, tt := range tests {
		t.Logf("running: %s", tt.name)

		userRepo := &mocks.UserRepo{}
		userService := service.NewUserDBService(&store.Store{User: userRepo})
		tt.mock(userRepo)

		got, err := userService.CreateUser(tt.input)
		if tt.wantErr {
			assert.Error(t, err)
			assert.Equal(t, tt.err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		userRepo.AssertExpectations(t)
	}
}

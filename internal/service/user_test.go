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
				userRepo.On("GetUserByUsername", "test").Return(&model.User{ID: 1, Username: "test", FullName: "test test", PhotoURL: "test.jpg"}, nil)
			},
			input: "test",
			want:  &model.User{ID: 1, Username: "test", FullName: "test test", PhotoURL: "test.jpg"},
		},
		{
			name: "Error [User not found]",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUserByUsername", "test").Return(nil, nil)
			},
			input: "test",
			want:  nil,
		},
		{
			name: "Error: [Store error]",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUserByUsername", "test").Return(nil, errors.New("error while getting user: some error"))
			},
			input:   "test",
			wantErr: true,
			err:     errors.New("[User] Service.GetUserByUsername error: error while getting user: some error"),
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
				userRepo.On("GetUserByUsername", input.Username).Return(nil, nil)
				userRepo.On("CreateUser", input).Return(1, nil)
			},
			input: input,
			want:  1,
		},
		{
			name: "Error: [User is exist]",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUserByUsername", input.Username).Return(input, nil)
			},
			input:   input,
			wantErr: true,
			err:     errors.New("user with username test is exist"),
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

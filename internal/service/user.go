package service

import (
	"fmt"

	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/internal/store"
)

type UserDBService struct {
	store *store.Store
}

func NewUserDBService(store *store.Store) *UserDBService {
	return &UserDBService{store: store}
}

func (s *UserDBService) GetUsers() ([]model.User, error) {
	users, err := s.store.User.GetUsers()
	if err != nil {
		return nil, fmt.Errorf("[User] Service.GetUser error: %w", err)
	}

	if users == nil {
		return nil, fmt.Errorf("user not found")
	}

	return users, nil
}

func (s *UserDBService) GetUserByUsername(username string) (*model.User, error) {
	user, err := s.store.User.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("[User] Service.GetUserByUsername error: %w", err)
	}
	if user == nil {
		return nil, nil
	}

	return user, nil
}

func (s *UserDBService) CreateUser(user *model.User) error {
	candidate, err := s.store.User.GetUserByUsername(user.Username)
	if err != nil {
		return err
	}

	if candidate != nil {
		return fmt.Errorf("User with username %s is exist", user.Username)
	}

	_, err = s.store.User.CreateUser(user)
	if err != nil {
		return fmt.Errorf("[User] Service.CreateUser error: %w", err)
	}

	return nil
}

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

func (s *UserDBService) CreateUser(user *model.User) (int, error) {
	if user.Username == "" {
		return 0, fmt.Errorf("username length should be more")
	}

	candidate, err := s.GetUserByUsername(user.Username)
	if err != nil {
		return candidate.ID, err
	}

	if candidate != nil {
		return candidate.ID, fmt.Errorf("user with username %s is exist", user.Username)
	}

	id, err := s.store.User.CreateUser(user)
	if err != nil {
		return id, fmt.Errorf("[User] Service.CreateUser error: %w", err)
	}

	return id, nil
}

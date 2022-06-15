package service

import (
	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/internal/store"
	"github.com/VladPetriv/tg_scanner/pkg/utils"
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
		return nil, &utils.ServiceError{
			ServiceName:       "User",
			ServiceMethodName: "GetUserByUsername",
			ErrorValue:        err,
		}
	}

	if user == nil {
		return nil, nil
	}

	return user, nil
}

func (s *UserDBService) CreateUser(user *model.User) (int, error) {
	utils.ValidateTelegramUser(user)

	candidate, err := s.GetUserByUsername(user.Username)
	if err != nil {
		return candidate.ID, err
	}

	if candidate != nil {
		return candidate.ID, &utils.RecordIsExistError{RecordName: "user", Name: user.Username}
	}

	id, err := s.store.User.CreateUser(user)
	if err != nil {
		return id, &utils.ServiceError{
			ServiceName:       "User",
			ServiceMethodName: "CreateUser",
			ErrorValue:        err,
		}
	}

	return id, nil
}

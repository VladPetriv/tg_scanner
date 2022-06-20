package service

import (
	"fmt"

	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/internal/store"
	"github.com/VladPetriv/tg_scanner/pkg/utils"
)

type ReplieDBService struct {
	store *store.Store
}

func NewReplieDBService(store *store.Store) *ReplieDBService {
	return &ReplieDBService{store: store}
}

func (s *ReplieDBService) GetReplieByName(name string) (*model.Replie, error) {
	replie, err := s.store.Replie.GetReplieByName(name)
	if err != nil {
		return nil, &utils.ServiceError{
			ServiceName:       "Replie",
			ServiceMethodName: "GetReplieByName",
			ErrorValue:        err,
		}
	}

	if replie == nil {
		return nil, nil
	}

	return replie, nil
}

func (s *ReplieDBService) CreateReplie(replie *model.Replie) error {
	if replie.Title == "" {
		return fmt.Errorf("title length should be more")
	}

	candidate, err := s.GetReplieByName(replie.Title)
	if err != nil {
		return err
	}

	fmt.Println(candidate)

	if candidate != nil && candidate.MessageID == replie.MessageID {
		return &utils.RecordIsExistError{RecordName: "replie", Name: replie.Title}
	}

	_, err = s.store.Replie.CreateReplie(replie)
	if err != nil {
		return &utils.ServiceError{
			ServiceName:       "Replie",
			ServiceMethodName: "CreateReplie",
			ErrorValue:        err,
		}
	}

	return nil
}

package service

import (
	"fmt"

	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/internal/store"
)

type ReplieDBService struct {
	store *store.Store
}

func NewReplieDBService(store *store.Store) *ReplieDBService {
	return &ReplieDBService{store: store}
}

func (s *ReplieDBService) GetReplie(replieID int) (*model.Replie, error) {
	replie, err := s.store.Replie.GetReplie(replieID)
	if err != nil {
		return nil, fmt.Errorf("[Replie] Service.GetReplie error: %w", err)
	}

	if replie == nil {
		return nil, fmt.Errorf("replie not found")
	}

	return replie, nil
}

func (s *ReplieDBService) GetReplieByName(name string) (*model.Replie, error) {
	replie, err := s.store.Replie.GetReplieByName(name)
	if err != nil {
		return nil, fmt.Errorf("[Replie] Service.GetReplieByName error: %w", err)
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

	if candidate != nil {
		return fmt.Errorf("replie with name %s is exist", replie.Title)
	}

	_, err = s.store.Replie.CreateReplie(replie)
	if err != nil {
		return fmt.Errorf("[Replie] Service.CreateReplie error: %w", err)
	}

	return nil
}

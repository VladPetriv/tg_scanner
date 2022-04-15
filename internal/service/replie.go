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

func (s *ReplieDBService) GetReplies() ([]model.Replie, error) {
	replies, err := s.store.Replie.GetReplies()
	if err != nil {
		return nil, fmt.Errorf("[Replie] Service.GetReplies error: %w", err)
	}

	if replies == nil {
		return nil, fmt.Errorf("replies not found")
	}

	return replies, nil
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
	_, err := s.store.Replie.CreateReplie(replie)
	if err != nil {
		return fmt.Errorf("[Replie] Service.CreateReplie error: %w", err)
	}

	return nil
}

func (s *ReplieDBService) DeleteReplie(replieID int) error {
	err := s.store.Replie.DeleteReplie(replieID)
	if err != nil {
		return fmt.Errorf("[Replie] Service.DeleteReplie error: %w", err)
	}

	return nil
}

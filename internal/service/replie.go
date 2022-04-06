package service

import (
	"fmt"

	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/internal/store"
)

type ReplieDbService struct {
	store *store.Store
}

func NewReplieDbService(store *store.Store) *ReplieDbService {
	return &ReplieDbService{store: store}
}

func (s *ReplieDbService) GetReplies() ([]model.Replie, error) {
	replies, err := s.store.Replie.GetReplies()
	if err != nil {
		return nil, fmt.Errorf("[Replie] Service.GetReplies error: %w", err)
	}

	if replies == nil {
		return nil, fmt.Errorf("replies not found")
	}

	return replies, nil
}

func (s *ReplieDbService) GetReplie(replieId int) (*model.Replie, error) {
	replie, err := s.store.Replie.GetReplie(replieId)
	if err != nil {
		return nil, fmt.Errorf("[Replie] Service.GetReplie error: %w", err)
	}

	if replie == nil {
		return nil, fmt.Errorf("replie not found")
	}

	return replie, nil
}

func (s *ReplieDbService) GetReplieByName(name string) (*model.Replie, error) {
	replie, err := s.store.Replie.GetReplieByName(name)
	if err != nil {
		return nil, fmt.Errorf("[Replie] Service.GetReplieByName error: %w", err)
	}

	if replie == nil {
		return nil, nil
	}

	return replie, nil
}

func (s *ReplieDbService) CreateReplie(replie *model.Replie) error {
	_, err := s.store.Replie.CreateReplie(replie)
	if err != nil {
		return fmt.Errorf("[Replie] Service.CreateReplie error: %v", err)
	}

	return nil
}

func (s *ReplieDbService) DeleteReplie(replieId int) error {
	err := s.store.Replie.DeleteReplie(replieId)
	if err != nil {
		return fmt.Errorf("[Replie] Service.DeleteReplie error: %w", err)
	}

	return nil
}

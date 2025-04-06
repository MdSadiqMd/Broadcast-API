package services

import (
	"github.com/MdSadiqMd/Broadcast-API/internal/models"
	"github.com/MdSadiqMd/Broadcast-API/internal/repositories"
	"gorm.io/gorm"
)

type BroadcastService struct {
	repo repositories.BroadcastRepository
}

func NewBroadcastService(db *gorm.DB) *BroadcastService {
	return &BroadcastService{
		repo: *repositories.NewBroadcastRepository(db),
	}
}

func (s *BroadcastService) CreateBroadcast(broadcast *models.Broadcast) (*models.Broadcast, error) {
	createdBroadcast, err := s.repo.CreateBroadcast(broadcast)
	if err != nil {
		return nil, err
	}
	return &createdBroadcast, nil
}

func (s *BroadcastService) GetBroadcastByID(id uint) (*models.Broadcast, error) {
	broadcast, err := s.repo.GetBroadcastByID(id)
	if err != nil {
		return nil, err
	}
	return broadcast, nil
}

func (s *BroadcastService) UpdateBroadcast(broadcast *models.Broadcast) (*models.Broadcast, error) {
	updatedBroadcast, err := s.repo.UpdateBroadcast(broadcast)
	if err != nil {
		return nil, err
	}
	return updatedBroadcast, nil
}

func (s *BroadcastService) SendBroadcast(broadcast_id, scheduled_at string) (id string, error error) {
	broadcast, err := s.repo.SendBroadcast(broadcast_id, scheduled_at)
	if err != nil {
		return "", err
	}
	return broadcast, nil
}

func (s *BroadcastService) ListBroadcasts() ([]models.Broadcast, error) {
	broadcasts, err := s.repo.ListBroadcasts()
	if err != nil {
		return nil, err
	}
	return broadcasts, nil
}

func (s *BroadcastService) DeleteBroadcast(id uint) error {
	err := s.repo.DeleteBroadcast(id)
	if err != nil {
		return err
	}
	return nil
}

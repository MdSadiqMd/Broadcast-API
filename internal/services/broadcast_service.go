package services

import (
	"time"

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

func (s *BroadcastService) UpdateBroadcast(id uint, broadcast *models.Broadcast) (*models.Broadcast, error) {
	existingBoradcast, err := s.repo.GetBroadcastByID(id)
	if err != nil {
		return nil, err
	}

	existingBoradcast.Name = broadcast.Name
	existingBoradcast.AudienceID = broadcast.AudienceID
	existingBoradcast.CampaignID = broadcast.CampaignID
	existingBoradcast.UserID = broadcast.UserID
	existingBoradcast.From = broadcast.From
	existingBoradcast.Subject = broadcast.Subject
	existingBoradcast.ReplyTo = broadcast.ReplyTo
	existingBoradcast.HTML = broadcast.HTML
	existingBoradcast.Text = broadcast.Text
	existingBoradcast.Status = broadcast.Status
	existingBoradcast.SentAt = broadcast.SentAt
	existingBoradcast.Campaign = broadcast.Campaign
	existingBoradcast.User = broadcast.User
	existingBoradcast.ScheduledAt = broadcast.ScheduledAt
	existingBoradcast.UpdatedAt = time.Now()

	updatedBroadcast, err := s.repo.UpdateBroadcast(existingBoradcast)
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

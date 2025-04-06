package repositories

import (
	"github.com/MdSadiqMd/Broadcast-API/internal/models"
	"gorm.io/gorm"
)

type BroadcastRepository struct {
	db *gorm.DB
}

func NewBroadcastRepository(db *gorm.DB) *BroadcastRepository {
	return &BroadcastRepository{
		db: db,
	}
}

func (r *BroadcastRepository) CreateBroadcast(broadcast *models.Broadcast) (models.Broadcast, error) {
	err := r.db.Create(broadcast).Error
	return *broadcast, err
}

func (r *BroadcastRepository) GetBroadcastByID(id uint) (*models.Broadcast, error) {
	var broadcast models.Broadcast
	err := r.db.Find(&broadcast, id).Error
	if err != nil {
		return nil, err
	}
	return &broadcast, nil
}

func (r *BroadcastRepository) UpdateBroadcast(broadcast *models.Broadcast) (*models.Broadcast, error) {
	err := r.db.Save(broadcast).Error
	return broadcast, err
}

func (r *BroadcastRepository) SendBroadcast(broadcast_id, scheduled_at string) (id string, error error) {
	err := r.db.Table("broadcasts").Where("id = ?", broadcast_id).Update("scheduled_at", scheduled_at).Error
	return broadcast_id, err
}

func (r *BroadcastRepository) ListBroadcasts() ([]models.Broadcast, error) {
	var broadcasts []models.Broadcast
	err := r.db.Find(&broadcasts).Error
	if err != nil {
		return nil, err
	}
	return broadcasts, nil
}

func (r *BroadcastRepository) DeleteBroadcast(id uint) error {
	return r.db.Delete(&models.Broadcast{}, id).Error
}

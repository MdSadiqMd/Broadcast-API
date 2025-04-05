package repositories

import (
	"github.com/MdSadiqMd/Broadcast-API/internal/models"
	"gorm.io/gorm"
)

type CampaignRepository struct {
	db *gorm.DB
}

func NewCampaignRepository(db *gorm.DB) *CampaignRepository {
	return &CampaignRepository{
		db: db,
	}
}

func (r *CampaignRepository) Create(campaign *models.Campaign) error {
	return r.db.Create(campaign).Error
}

func (r *CampaignRepository) GetAllCampaigns() ([]models.Campaign, error) {
	var campaigns []models.Campaign
	err := r.db.Find(&campaigns).Error
	if err != nil {
		return nil, err
	}
	return campaigns, nil
}

func (r *CampaignRepository) GetCampaignByID(id uint) (*models.Campaign, error) {
	var campaign models.Campaign
	err := r.db.Find(&campaign, id).Error
	if err != nil {
		return nil, err
	}
	return &campaign, nil
}

func (r *CampaignRepository) DeleteCampaign(id uint) error {
	return r.db.Delete(&models.Campaign{}, id).Error
}

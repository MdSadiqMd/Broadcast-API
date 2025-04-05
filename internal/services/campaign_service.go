package services

import (
	"github.com/MdSadiqMd/Broadcast-API/internal/models"
	"github.com/MdSadiqMd/Broadcast-API/internal/repositories"
	"gorm.io/gorm"
)

type CampaignService struct {
	repo repositories.CampaignRepository
}

func NewCampaignService(db *gorm.DB) *CampaignService {
	return &CampaignService{
		repo: *repositories.NewCampaignRepository(db),
	}
}

func (s *CampaignService) CreateCampaign(campaign *models.Campaign) (*models.Campaign, error) {
	createdCampaign, err := s.repo.Create(campaign)
	if err != nil {
		return nil, err
	}
	return &createdCampaign, nil
}

func (s *CampaignService) GetAllCampaigns() ([]models.Campaign, error) {
	compaigns, err := s.repo.GetAllCampaigns()
	if err != nil {
		return nil, err
	}
	return compaigns, nil
}

func (s *CampaignService) GetCampaignByID(id uint) (*models.Campaign, error) {
	compaign, err := s.repo.GetCampaignByID(id)
	if err != nil {
		return nil, err
	}
	return compaign, nil
}

func (s *CampaignService) DeleteCampaign(id uint) error {
	err := s.repo.DeleteCampaign(id)
	if err != nil {
		return err
	}
	return nil
}

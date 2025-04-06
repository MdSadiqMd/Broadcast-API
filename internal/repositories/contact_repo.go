package repositories

import (
	"github.com/MdSadiqMd/Broadcast-API/internal/models"
	"gorm.io/gorm"
)

type ContactRepository struct {
	db *gorm.DB
}

func NewContactRepository(db *gorm.DB) *ContactRepository {
	return &ContactRepository{
		db: db,
	}
}

func (r *ContactRepository) CreateContact(contact *models.Contact) (models.Contact, error) {
	err := r.db.Create(contact).Error
	return *contact, err
}

func (r *ContactRepository) GetAllContacts(campaignID uint) ([]models.Contact, error) {
	var contacts []models.Contact
	result := r.db.Joins("JOIN campaign_audiences ON campaign_audiences.contact_id = contacts.id").
		Where("campaign_audiences.campaign_id = ?", campaignID).
		Find(&contacts)

	if result.Error != nil {
		return nil, result.Error
	}

	return contacts, nil
}

func (r *ContactRepository) GetContactByID(id uint) (*models.Contact, error) {
	var contact models.Contact
	err := r.db.Find(&contact, id).Error
	if err != nil {
		return nil, err
	}
	return &contact, nil
}

func (r *ContactRepository) UpdateContact(contact *models.Contact) (models.Contact, error) {
	err := r.db.Save(contact).Error
	return *contact, err
}

func (r *ContactRepository) DeleteContact(id uint) error {
	return r.db.Delete(&models.Contact{}, id).Error
}

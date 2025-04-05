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

func (r *ContactRepository) GetAllContacts(audienceId uint) ([]models.Contact, error) {
	var contacts []models.Contact
	err := r.db.Where("audience_id = ?", audienceId).Find(&contacts).Error
	return contacts, err
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

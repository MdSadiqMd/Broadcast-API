package services

import (
	"github.com/MdSadiqMd/Broadcast-API/internal/models"
	"github.com/MdSadiqMd/Broadcast-API/internal/repositories"
	"gorm.io/gorm"
)

type ContactService struct {
	repo *repositories.ContactRepository
}

func NewContactService(db *gorm.DB) *ContactService {
	return &ContactService{
		repo: repositories.NewContactRepository(db),
	}
}

func (s *ContactService) CreateContact(contact *models.Contact) (*models.Contact, error) {
	createdContact, err := s.repo.CreateContact(contact)
	if err != nil {
		return nil, err
	}
	return &createdContact, nil
}

func (s *ContactService) GetAllContacts(audienceId uint) ([]models.Contact, error) {
	contacts, err := s.repo.GetAllContacts(audienceId)
	if err != nil {
		return nil, err
	}
	return contacts, nil
}

func (s *ContactService) GetContactByID(id uint) (*models.Contact, error) {
	contact, err := s.repo.GetContactByID(id)
	if err != nil {
		return nil, err
	}
	return contact, nil
}

func (s *ContactService) UpdateContact(id uint, contact *models.Contact) (*models.Contact, error) {
	existingContact, err := s.repo.GetContactByID(id)
	if err != nil {
		return nil, err
	}

	existingContact.FirstName = contact.FirstName
	existingContact.LastName = contact.LastName
	existingContact.Email = contact.Email

	updatedContact, err := s.repo.UpdateContact(existingContact)
	if err != nil {
		return nil, err
	}

	return &updatedContact, nil
}

func (s *ContactService) DeleteContact(id uint) error {
	err := s.repo.DeleteContact(id)
	if err != nil {
		return err
	}
	return nil
}

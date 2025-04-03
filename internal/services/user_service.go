package services

import (
	"errors"
	"time"

	"github.com/MdSadiqMd/Broadcast-API/internal/models"
	"github.com/MdSadiqMd/Broadcast-API/internal/repositories"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		repo: repositories.NewUserRepository(db),
	}
}

func (s *UserService) CreateUser(username, password, email, role string) (*models.User, error) {
	existingUser, _ := s.repo.FindByUsernameOrEmail(username, email)
	if existingUser != nil {
		return nil, errors.New("username or email already exists")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:  username,
		Password:  string(hashPassword),
		Email:     email,
		Role:      role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Authenticate(username, password string) (*models.User, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *UserService) UserExists(id uint) (bool, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return user != nil, nil
}

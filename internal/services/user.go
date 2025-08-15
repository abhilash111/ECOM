package services

import (
	"github.com/abhilash111/ecom/internal/models"
	"github.com/abhilash111/ecom/internal/repository"
)

type UserService interface {
	GetUserByID(id uint) (*models.User, error)
	GetCurrentSubscription(userID uint) (*models.UserSubscription, error)
	UpdateUser(user *models.User) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.GetUserByID(id)
}

func (s *userService) GetCurrentSubscription(userID uint) (*models.UserSubscription, error) {
	return s.userRepo.GetCurrentSubscription(userID)
}

func (s *userService) UpdateUser(user *models.User) error {
	return s.userRepo.UpdateUser(user)
}

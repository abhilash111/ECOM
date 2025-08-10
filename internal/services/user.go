package services

import (
	"github.com/abhilash111/ecom/internal/models"
	"github.com/abhilash111/ecom/internal/repository"
)

type UserService interface {
	GetUserByID(id uint) (*models.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

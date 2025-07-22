package services

import (
	"github.com/abhilash111/ecom/internal/models"
	"github.com/abhilash111/ecom/internal/repository"
)

type UserService interface {
	CreateUser(email, phoneNumber, cognitoID string, roles []models.Role) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByPhone(phoneNumber string) (*models.User, error)
	GetUserByCognitoID(cognitoID string) (*models.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) CreateUser(email, phoneNumber, cognitoID string, roles []models.Role) (*models.User, error) {
	return s.userRepo.CreateUser(email, phoneNumber, cognitoID, roles)
}

func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	return s.userRepo.GetUserByEmail(email)
}

func (s *userService) GetUserByPhone(phoneNumber string) (*models.User, error) {
	return s.userRepo.GetUserByPhone(phoneNumber)
}

func (s *userService) GetUserByCognitoID(cognitoID string) (*models.User, error) {
	return s.userRepo.GetUserByCognitoID(cognitoID)
}

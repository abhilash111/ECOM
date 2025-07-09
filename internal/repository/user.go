package repository

import (
	"github.com/abhilash111/ecom/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(email, phoneNumber, cognitoID string, roles []models.Role) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByPhone(phoneNumber string) (*models.User, error)
	GetUserByCognitoID(cognitoID string) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(email, phoneNumber, cognitoID string, roles []models.Role) (*models.User, error) {
	// Convert roles to comma-separated string
	var rolesStr string
	for i, role := range roles {
		if i > 0 {
			rolesStr += ","
		}
		rolesStr += string(role)
	}

	user := &models.User{
		Email:       email,
		PhoneNumber: phoneNumber,
		CognitoID:   cognitoID,
		Roles:       rolesStr,
	}

	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByPhone(phoneNumber string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("phone_number = ?", phoneNumber).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByCognitoID(cognitoID string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("cognito_id = ?", cognitoID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

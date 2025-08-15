package repository

import (
	"github.com/abhilash111/ecom/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	FindUserByPhone(phone string) (*models.User, error)
	FindUserByEmail(email string) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
	UpdateUser(user *models.User) error
	CreateSubscription(subscription *models.UserSubscription) error
	GetCurrentSubscription(userID uint) (*models.UserSubscription, error)
	GetPackageByType(packageType models.PackageType) (*models.Package, error)
	CreateRefreshToken(token *models.RefreshToken) error
	GetRefreshToken(token string) (*models.RefreshToken, error)
	RevokeRefreshToken(token string) error
	CreateSession(session *models.Session) error
	GetSession(sessionID string) (*models.Session, error)
	DeleteSession(sessionID string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindUserByPhone(phone string) (*models.User, error) {
	var user models.User
	err := r.db.Where("phone_number = ?", phone).First(&user).Error
	return &user, err
}

func (r *userRepository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Subscriptions").Preload("Subscriptions.Package").First(&user, id).Error
	return &user, err
}

func (r *userRepository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) CreateSubscription(subscription *models.UserSubscription) error {
	return r.db.Create(subscription).Error
}

func (r *userRepository) GetCurrentSubscription(userID uint) (*models.UserSubscription, error) {
	var subscription models.UserSubscription
	err := r.db.Where("user_id = ? AND expires_at > NOW()", userID).
		Order("expires_at desc").
		First(&subscription).
		Error
	return &subscription, err
}

func (r *userRepository) GetPackageByType(packageType models.PackageType) (*models.Package, error) {
	var pkg models.Package
	err := r.db.Where("name = ?", packageType).First(&pkg).Error
	return &pkg, err
}

func (r *userRepository) CreateRefreshToken(token *models.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *userRepository) GetRefreshToken(token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	err := r.db.Where("token = ?", token).First(&refreshToken).Error
	return &refreshToken, err
}

func (r *userRepository) RevokeRefreshToken(token string) error {
	return r.db.Model(&models.RefreshToken{}).
		Where("token = ?", token).
		Update("revoked", true).Error
}

func (r *userRepository) CreateSession(session *models.Session) error {
	return r.db.Create(session).Error
}

func (r *userRepository) GetSession(sessionID string) (*models.Session, error) {
	var session models.Session
	err := r.db.Where("session_id = ?", sessionID).First(&session).Error
	return &session, err
}

func (r *userRepository) DeleteSession(sessionID string) error {
	return r.db.Where("session_id = ?", sessionID).Delete(&models.Session{}).Error
}

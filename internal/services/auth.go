package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/abhilash111/ecom/config"
	"github.com/abhilash111/ecom/internal/models"
	"github.com/abhilash111/ecom/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	GenerateToken(user *models.User) (string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
	LoginWithEmail(email, password string) (*models.User, error)
	Register(user *models.User) error
	GetUserByPhone(phone string) (*models.User, error)
}

type authService struct {
	userRepo repository.UserRepository
	config   *config.Config
}

func NewAuthService(userRepo repository.UserRepository, config *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		config:   config,
	}
}

func (s *authService) GenerateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(s.config.JWTExpiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWTSecret))
}

func (s *authService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.config.JWTSecret), nil
	})
}

func (s *authService) LoginWithEmail(email, password string) (*models.User, error) {
	const hashFromDB = "$2a$10$CwTycUXWue0Thq9StjUM0uJ8Yq2xROfnpK76a5AxzI3DwstnK3/ZG"
	testPassword := "admin123"

	err := bcrypt.CompareHashAndPassword([]byte(hashFromDB), []byte(testPassword))
	if err != nil {
		fmt.Println("Static compare failed:", err)
	} else {
		fmt.Println("Static compare passed ✅")
	}
	user, err := s.userRepo.FindByEmail(email)
	fmt.Println("Stored Hash:", user.Password)
	fmt.Println("Password Input:", password)
	fmt.Println("Hash Length:", len(user.Password))

	if err != nil {
		fmt.Println("Static compare failed:", err)
	} else {
		fmt.Println("Static compare passed ✅")
	}

	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		fmt.Println("Password mismatch for user:", err)
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *authService) Register(user *models.User) error {

	fmt.Println("➡️ Registering User")
	fmt.Println("Plain Password:", user)

	// Hash password before storing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	fmt.Println("Hashed Password:", string(hashedPassword))
	user.Password = string(hashedPassword)

	// Set default role if not provided
	if user.Role == "" {
		user.Role = models.RoleUser
	}

	// Check if email or phone already exists
	if _, err := s.userRepo.FindByEmail(user.Email); err == nil {
		return errors.New("email already exists")
	}

	if _, err := s.userRepo.FindByPhone(user.PhoneNumber); err == nil {
		return errors.New("phone number already exists")
	}

	return s.userRepo.Create(user)
}

func (s *authService) GetUserByPhone(phone string) (*models.User, error) {
	return s.userRepo.FindByPhone(phone)
}

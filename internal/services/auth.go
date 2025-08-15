package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/abhilash111/ecom/internal/repository"

	"github.com/abhilash111/ecom/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	RegisterUser(request models.RegisterRequest) (*models.User, error)
	LoginWithEmail(email, password string) (*models.User, error)
	GenerateTokenPair(user *models.User) (*models.TokenPair, error)
	FindUserByPhone(phone string) (*models.User, error) // Added this method
	RefreshAccessToken(refreshToken string) (*models.TokenPair, error)
	RevokeRefreshToken(token string) error
	CreateSession(user *models.User, userAgent, ipAddress string) (*models.Session, error)
	ValidateSession(sessionID string) (bool, *models.User, error)
	DeleteSession(sessionID string) error
	ParseJWT(tokenString string) (*jwt.Token, error)
}

type authService struct {
	userRepo   repository.UserRepository
	redisRepo  repository.RedisRepository
	jwtSecret  string
	accessExp  time.Duration
	refreshExp time.Duration
}

func NewAuthService(
	userRepo repository.UserRepository,
	redisRepo repository.RedisRepository,
	jwtSecret string,
	accessExp time.Duration,
	refreshExp time.Duration,
) AuthService {
	return &authService{
		userRepo:   userRepo,
		redisRepo:  redisRepo,
		jwtSecret:  jwtSecret,
		accessExp:  accessExp,
		refreshExp: refreshExp,
	}
}

func (s *authService) RegisterUser(request models.RegisterRequest) (*models.User, error) {
	_, err := s.userRepo.FindUserByPhone(request.PhoneNumber)
	if err == nil {
		return nil, errors.New("user with this phone number already exists")
	}

	_, err = s.userRepo.FindUserByEmail(request.Email)
	if err == nil {
		return nil, errors.New("user with this email already exists")
	}

	pkg, err := s.userRepo.GetPackageByType(request.SubscriptionPack)
	if err != nil {
		return nil, errors.New("invalid subscription package")
	}

	user := &models.User{
		Name:        request.Name,
		PhoneNumber: request.PhoneNumber,
		Email:       request.Email,
		IsPassword:  request.Password != "",
		Role:        models.RoleUser,
	}

	if request.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedPassword)
	}

	err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	subscription := &models.UserSubscription{
		UserID:               user.ID,
		PackageID:            pkg.ID,
		StartsAt:             now,
		ExpiresAt:            now.AddDate(1, 0, 0),
		RemainingViewCredits: pkg.ViewListingLimit,
		RemainingAddCredits:  pkg.AddListingLimit,
	}

	err = s.userRepo.CreateSubscription(subscription)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) LoginWithEmail(email, password string) (*models.User, error) {
	user, err := s.userRepo.FindUserByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !user.IsPassword {
		return nil, errors.New("password login not enabled for this user")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}

func (s *authService) generateAccessToken(user *models.User) (string, error) {
	expiryTime := time.Now().Add(s.accessExp)

	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   expiryTime.Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

func (s *authService) GenerateTokenPair(user *models.User) (*models.TokenPair, error) {
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken := uuid.New().String()
	refreshTokenRecord := models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(s.refreshExp),
	}

	if err := s.userRepo.CreateRefreshToken(&refreshTokenRecord); err != nil {
		return nil, err
	}

	return &models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) RefreshAccessToken(refreshToken string) (*models.TokenPair, error) {
	tokenRecord, err := s.userRepo.GetRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if tokenRecord.Revoked || tokenRecord.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("refresh token expired or revoked")
	}

	user, err := s.userRepo.GetUserByID(tokenRecord.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return s.GenerateTokenPair(user)
}

func (s *authService) RevokeRefreshToken(token string) error {
	return s.userRepo.RevokeRefreshToken(token)
}

func (s *authService) CreateSession(user *models.User, userAgent, ipAddress string) (*models.Session, error) {
	session := &models.Session{
		UserID:    user.ID,
		SessionID: uuid.New().String(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
		UserAgent: userAgent,
		IPAddress: ipAddress,
	}

	if err := s.userRepo.CreateSession(session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *authService) ValidateSession(sessionID string) (bool, *models.User, error) {
	fmt.Print("Validating session with ID:", sessionID, "\n")
	session, err := s.userRepo.GetSession(sessionID)
	if err != nil {
		return false, nil, err
	}

	if session.ExpiresAt.Before(time.Now()) {
		return false, nil, nil
	}

	user, err := s.userRepo.GetUserByID(session.UserID)
	if err != nil {
		return false, nil, err
	}

	return true, user, nil
}

func (s *authService) DeleteSession(sessionID string) error {
	return s.userRepo.DeleteSession(sessionID)
}

func (s *authService) ParseJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})
}

func (s *authService) FindUserByPhone(phone string) (*models.User, error) {
	return s.userRepo.FindUserByPhone(phone)
}

package services

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/abhilash111/ecom/config"
	"github.com/abhilash111/ecom/internal/models"
	"github.com/abhilash111/ecom/internal/repository"
)

type OTPService interface {
	GenerateOTP(phoneNumber string) (*models.OTP, error)
	VerifyOTP(phoneNumber, code string) (bool, error)
}

type otpService struct {
	redisRepo repository.RedisRepository
	config    *config.Config
}

func NewOTPService(redisRepo repository.RedisRepository, config *config.Config) OTPService {
	return &otpService{
		redisRepo: redisRepo,
		config:    config,
	}
}

func (s *otpService) GenerateOTP(phoneNumber string) (*models.OTP, error) {
	// Generate 6-digit OTP
	code, err := generateRandomNumber(6)
	if err != nil {
		return nil, err
	}

	otp := &models.OTP{
		PhoneNumber: phoneNumber,
		Code:        code,
		ExpiresAt:   int64(s.config.OTPExpiration.Seconds()),
	}

	err = s.redisRepo.StoreOTP(otp)
	if err != nil {
		return nil, err
	}

	return otp, nil
}

func (s *otpService) VerifyOTP(phoneNumber, code string) (bool, error) {
	storedOTP, err := s.redisRepo.GetOTP(phoneNumber)
	fmt.Println("Retrieved OTP:", storedOTP)
	if err != nil {
		fmt.Println("Error retrieving OTP:", err)
		return false, err
	}

	fmt.Println("Verifying OTP:", storedOTP.Code, "against provided code:", code)
	if storedOTP.Code != code {
		return false, nil
	}

	// Delete OTP after verification
	err = s.redisRepo.DeleteOTP(phoneNumber)
	if err != nil {
		return false, err
	}

	return true, nil
}

func generateRandomNumber(length int) (string, error) {
	const digits = "0123456789"
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		result[i] = digits[num.Int64()]
	}
	return string(result), nil
}

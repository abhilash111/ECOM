package services

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/abhilash111/ecom/config"
	"github.com/abhilash111/ecom/internal/repository"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

const (
	otpLength     = 6
	otpExpiration = 5 * time.Minute
)

type OTPService interface {
	GenerateOTP() (string, error)
	SendOTP(phoneNumber string) (string, error)
	VerifyOTP(phoneNumber, otp string) (bool, error)
}

type otpService struct{}

func NewOTPService() OTPService {
	return &otpService{}
}

func (s *otpService) GenerateOTP() (string, error) {
	const digits = "0123456789"
	otp := make([]byte, otpLength)

	for i := range otp {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		otp[i] = digits[num.Int64()]
	}

	return string(otp), nil
}

func (s *otpService) SendOTP(phoneNumber string) (string, error) {
	otp, err := s.GenerateOTP()
	if err != nil {
		return "", err
	}

	// Store OTP in Redis
	err = repository.StoreOTP(phoneNumber, otp, otpExpiration)
	if err != nil {
		return "", err
	}

	// Invoke Lambda to send OTP
	sess := session.Must(session.NewSession())
	lambdaSvc := lambda.New(sess, aws.NewConfig().WithRegion(config.Envs.CognitoRegion))

	payload := fmt.Sprintf(`{"phoneNumber": "%s", "otp": "%s"}`, phoneNumber, otp)
	_, err = lambdaSvc.Invoke(&lambda.InvokeInput{
		FunctionName: aws.String(config.Envs.LambdaOTPARN),
		Payload:      []byte(payload),
	})

	if err != nil {
		// If Lambda fails, clean up the OTP
		_ = repository.DeleteOTP(phoneNumber)
		return "", err
	}

	return otp, nil
}

func (s *otpService) VerifyOTP(phoneNumber, otp string) (bool, error) {
	return repository.VerifyOTP(phoneNumber, otp)
}

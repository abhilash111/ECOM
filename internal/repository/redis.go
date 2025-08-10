package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/abhilash111/ecom/internal/models"
	"github.com/go-redis/redis/v8"
)

const (
	otpExpiration = 5 * time.Minute
)

type RedisRepository interface {
	StoreOTP(otp *models.OTP) error
	GetOTP(phoneNumber string) (*models.OTP, error)
	DeleteOTP(phoneNumber string) error
}

type redisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) RedisRepository {
	return &redisRepository{client: client}
}

func (r *redisRepository) StoreOTP(otp *models.OTP) error {
	ctx := context.Background()
	// Format the OTP Model and Print for debugging

	fmt.Println("Storing OTP:", otp.Code, "for phone number:", otp.PhoneNumber)
	return r.client.Set(ctx, otp.PhoneNumber, otp.Code, otpExpiration).Err()
}

func (r *redisRepository) GetOTP(phoneNumber string) (*models.OTP, error) {
	ctx := context.Background()
	fmt.Println("Retrieving OTP for phone number:", phoneNumber)
	code, err := r.client.Get(ctx, phoneNumber).Result()
	if err != nil {
		fmt.Println("Error getting Redis for OTP:", err)
		return nil, err
	}

	ttl, err := r.client.TTL(ctx, phoneNumber).Result()
	if err != nil {
		fmt.Println("Error getting TTL for OTP:", err)
		return nil, err
	}

	return &models.OTP{
		PhoneNumber: phoneNumber,
		Code:        code,
		ExpiresAt:   int64(ttl.Seconds()),
	}, nil
}

func (r *redisRepository) DeleteOTP(phoneNumber string) error {
	ctx := context.Background()
	return r.client.Del(ctx, phoneNumber).Err()
}

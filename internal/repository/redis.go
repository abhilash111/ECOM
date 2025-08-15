package repository

import (
	"context"
	"time"

	"github.com/abhilash111/ecom/internal/models"

	"github.com/go-redis/redis/v8"
)

type RedisRepository interface {
	StoreOTP(otp *models.OTP) error
	GetOTP(phoneNumber string) (*models.OTP, error)
	DeleteOTP(phoneNumber string) error
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
}

type redisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) RedisRepository {
	return &redisRepository{client: client}
}

func (r *redisRepository) StoreOTP(otp *models.OTP) error {
	ctx := context.Background()
	return r.client.Set(ctx, otp.PhoneNumber, otp.Code, time.Duration(otp.ExpiresAt)*time.Second).Err()
}

func (r *redisRepository) GetOTP(phoneNumber string) (*models.OTP, error) {
	ctx := context.Background()
	code, err := r.client.Get(ctx, phoneNumber).Result()
	if err != nil {
		return nil, err
	}
	return &models.OTP{
		PhoneNumber: phoneNumber,
		Code:        code,
	}, nil
}

func (r *redisRepository) DeleteOTP(phoneNumber string) error {
	ctx := context.Background()
	return r.client.Del(ctx, phoneNumber).Err()
}

func (r *redisRepository) Set(key string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *redisRepository) Get(key string) (string, error) {
	ctx := context.Background()
	return r.client.Get(ctx, key).Result()
}

func (r *redisRepository) Delete(key string) error {
	ctx := context.Background()
	return r.client.Del(ctx, key).Err()
}

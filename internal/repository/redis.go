package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/abhilash111/ecom/config"
	"github.com/go-redis/redis/v8"
)

var (
	redisClient *redis.Client
	ctx         = context.Background()
)

func InitRedis() error {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Envs.RedisHost, config.Envs.RedisPort),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := redisClient.Ping(ctx).Result()
	return err
}

func StoreOTP(phoneNumber, otp string, expiration time.Duration) error {
	err := redisClient.Set(ctx, fmt.Sprintf("otp:%s", phoneNumber), otp, expiration).Err()
	return err
}

func VerifyOTP(phoneNumber, otp string) (bool, error) {
	storedOTP, err := redisClient.Get(ctx, fmt.Sprintf("otp:%s", phoneNumber)).Result()
	if err == redis.Nil {
		return false, nil // OTP not found
	} else if err != nil {
		return false, err
	}

	return storedOTP == otp, nil
}

func DeleteOTP(phoneNumber string) error {
	_, err := redisClient.Del(ctx, fmt.Sprintf("otp:%s", phoneNumber)).Result()
	return err
}

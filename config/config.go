package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	RedisHost     string
	RedisPort     string
	RedisPassword string
	JWTSecret     string
	JWTExpiration time.Duration
	OTPExpiration time.Duration
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	jwtExp, err := time.ParseDuration(os.Getenv("JWT_EXPIRATION"))
	if err != nil {
		jwtExp = 24 * time.Hour // default 24 hours
	}

	otpExp, err := time.ParseDuration(os.Getenv("OTP_EXPIRATION") + "s")
	if err != nil {
		otpExp = 5 * time.Minute // default 5 minutes
	}

	return &Config{
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBUser:        os.Getenv("DB_USERNAME"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPort:     os.Getenv("REDIS_PORT"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		JWTExpiration: jwtExp,
		OTPExpiration: otpExp,
	}
}

package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost             string
	Port                   string
	DBUser                 string
	DBPassword             string
	DBAddress              string
	DBName                 string
	DBHost                 string // Added for clarity, though not used in the connection string
	DBPort                 string // Added for clarity, though not used in the connection string
	JWTExpirationInSeconds int64
	JWTSecret              string
	RedisHost              string
	RedisPort              string
	CognitoRegion          string
	CognitoPoolID          string
	CognitoAppID           string
	LambdaOTPARN           string
	CognitoAppSecret       string // Added for secret hash calculation
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()
	return Config{
		PublicHost: getEnv("PUBLIC_HOST", "http://localhost"),
		Port:       getEnv("PORT", "8080"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "kulkarni11"), // Also password_
		// DBAddress:  fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306")),
		DBHost:                 getEnv("DB_HOST", "127.0.0.1"),
		DBPort:                 getEnv("DB_PORT", "3307"),
		DBAddress:              fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306")), // use this for running locally
		DBName:                 getEnv("DB_NAME", "ecom"),
		JWTExpirationInSeconds: getEnvAsInt("JWT_EXPIRATION_IN_SECONDS", 3600*24*7),
		JWTSecret:              getEnv("JWT_SECRET", "not-secret-anymore"),
		RedisHost:              getEnv("REDIS_HOST", "localhost"),
		RedisPort:              getEnv("REDIS_PORT", "6379"),
		CognitoRegion:          getEnv("COGNITO_REGION", "ap-south-1"),
		CognitoPoolID:          getEnv("COGNITO_POOL_ID", "ap-south-1_MvmPmDxqe"),
		CognitoAppID:           getEnv("COGNITO_APP_ID", "1u0q5scnpli2bbs51gr17e8qm1"),
		LambdaOTPARN:           getEnv("LAMBDA_OTP_ARN", "arn:aws:lambda:ap-south-1:448270596903:function:cognito-custom-sms-sender-dev-sendSMS"),
		CognitoAppSecret:       getEnv("COGNITO_APP_SECRET", "1amhese5o7uje9un5ug9or8fotjd6emen198fhfjfcfcg88il4qf"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		return i
	}

	return fallback
}

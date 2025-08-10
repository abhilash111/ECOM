package database

// Init Redis connection
import (
	"context"
	"fmt"
	"log"

	"github.com/abhilash111/ecom/config"
	"github.com/go-redis/redis/v8"
)

// ConnectRedis initializes the Redis connection
func ConnectRedis() *redis.Client {
	cfg := config.LoadConfig()
	redisAddr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)

	log.Println("Connecting to Redis at", redisAddr)
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Ping Redis to check connection
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	return client
}

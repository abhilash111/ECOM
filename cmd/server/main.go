package main

import (
	"log"

	"github.com/abhilash111/ecom/config"
	"github.com/abhilash111/ecom/internal/database"
	"github.com/abhilash111/ecom/internal/repository"
	"github.com/abhilash111/ecom/internal/routes"
	"github.com/abhilash111/ecom/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database connection
	db := database.ConnectDB()

	redisClient := database.ConnectRedis()

	userRepo := repository.NewUserRepository(db)
	redisRepo := repository.NewRedisRepository(redisClient)

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg)
	otpService := services.NewOTPService(redisRepo, cfg)
	userService := services.NewUserService(userRepo) // Add this line

	// Initialize Gin
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router, authService, otpService, userService)

	// Start server
	log.Println("Server starting on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}

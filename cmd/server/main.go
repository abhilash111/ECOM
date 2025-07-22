package main

import (
	"context"
	"fmt"
	"log"

	"github.com/abhilash111/ecom/config"
	"github.com/abhilash111/ecom/internal/controllers"
	"github.com/abhilash111/ecom/internal/models"
	"github.com/abhilash111/ecom/internal/repository"
	"github.com/abhilash111/ecom/internal/routes"
	"github.com/abhilash111/ecom/internal/services"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func LoadAWSConfig() aws.Config {
	cfg, err := awsConfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Unable to load AWS SDK config: %v", err)
	}
	return cfg
}

func main() {
	cfg := LoadAWSConfig()
	fmt.Println("AWS Region:", cfg.Region, cfg.Credentials)

	dsn := config.Envs.DBUser + ":" + config.Envs.DBPassword + "@tcp(" + config.Envs.DBHost + ":" + config.Envs.DBPort + ")/" + config.Envs.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize Redis
	err = repository.InitRedis()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Auto-migrate models
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize services
	userRepo := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	otpService := services.NewOTPService()
	authService := services.NewAuthService(userService, otpService)

	// Initialize controllers
	authController := controllers.NewAuthController(authService)
	userController := controllers.NewUserController(userService)

	// Create Gin router
	router := gin.Default()

	// Set up routes
	routes.SetupRoutes(router, authController, userController)

	// Start server
	log.Println("Server starting on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

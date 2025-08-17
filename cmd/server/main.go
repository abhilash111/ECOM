package main

import (
	"log"
	"time"

	"github.com/abhilash111/ecom/config"
	"github.com/abhilash111/ecom/internal/controllers"
	"github.com/abhilash111/ecom/internal/database"
	"github.com/abhilash111/ecom/internal/middleware"
	"github.com/abhilash111/ecom/internal/repository"
	"github.com/abhilash111/ecom/internal/routes"
	"github.com/abhilash111/ecom/internal/services"
	"github.com/gin-contrib/cors"
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
	authService := services.NewAuthService(
		userRepo,
		redisRepo,
		cfg.JWTSecret,
		15*time.Minute, // Access token expiry
		7*24*time.Hour, // Refresh token expiry
	)
	otpService := services.NewOTPService(redisRepo, cfg)
	userService := services.NewUserService(userRepo)

	// Initialize controllers
	authController := controllers.NewAuthController(authService, otpService)
	userController := controllers.NewUserController(userService)

	// Initialize middleware
	authMiddleware := middleware.AuthMiddleware(authService)

	// Setup router
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Your React app's origin
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.SetupRoutes(router, authController, userController, authMiddleware)

	// Start server
	log.Println("Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}

}

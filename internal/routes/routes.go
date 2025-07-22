package routes

import (
	"github.com/abhilash111/ecom/internal/controllers"
	"github.com/abhilash111/ecom/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, authController *controllers.AuthController, userController *controllers.UserController) {
	// Public routes
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/signup", authController.SignUp)
		authGroup.POST("/login", authController.Login)
		authGroup.POST("/phone-login/initiate", authController.InitiatePhoneLogin)
		authGroup.POST("/phone-login/verify", authController.VerifyPhoneLogin)
		authGroup.POST("/refresh", authController.RefreshToken)
	}

	// Protected routes
	protectedGroup := router.Group("/api")
	protectedGroup.Use(middleware.AuthMiddleware())
	{
		// User routes
		userGroup := protectedGroup.Group("/users")
		{
			userGroup.GET("/me", userController.GetCurrentUser)

			// Admin-only routes
			adminGroup := userGroup.Group("")
			adminGroup.Use(middleware.RoleMiddleware("admin"))
			{
				adminGroup.GET("", userController.GetAllUsers)
			}
		}
	}
}

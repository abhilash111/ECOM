package routes

import (
	"github.com/abhilash111/ecom/internal/controllers"
	"github.com/abhilash111/ecom/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	router *gin.Engine,
	authController *controllers.AuthController,
	userController *controllers.UserController,
	authMiddleware gin.HandlerFunc,
) {
	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login/email", authController.LoginWithEmail)
			auth.POST("/otp/request", authController.RequestOTP)
			auth.POST("/otp/verify", authController.LoginWithOTP)
			auth.POST("/logout", authController.Logout)
			auth.POST("/refresh", authController.RefreshToken)
		}

		user := api.Group("/user")
		user.Use(authMiddleware)
		{
			user.GET("/profile", userController.GetProfile)
		}

		admin := api.Group("/admin")
		admin.Use(authMiddleware, middleware.RoleMiddleware("admin"))
		{
			// Admin routes
		}
	}
}

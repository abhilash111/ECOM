package routes

import (
	"github.com/abhilash111/ecom/internal/controllers"
	"github.com/abhilash111/ecom/internal/middleware"
	"github.com/abhilash111/ecom/internal/services"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	router *gin.Engine,
	authService services.AuthService,
	otpService services.OTPService,
	userService services.UserService,
) {
	authController := controllers.NewAuthController(authService, otpService)
	userController := controllers.NewUserController(userService)

	// Public routes
	public := router.Group("/api/v1")
	{
		public.POST("/register", authController.Register)
		public.POST("/login/email", authController.LoginWithEmail)
		public.POST("/otp/request", authController.RequestOTP)
		public.POST("/login/otp", authController.LoginWithOTP)
	}

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(authService))
	{
		protected.GET("/profile", userController.GetProfile)
	}

	// Admin routes
	admin := router.Group("/api/v1/admin")
	admin.Use(middleware.AuthMiddleware(authService))
	admin.Use(middleware.RoleMiddleware("admin"))
	{
		// Add admin-specific routes here
	}
}

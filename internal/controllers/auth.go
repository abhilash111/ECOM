package controllers

import (
	"fmt"
	"net/http"

	"github.com/abhilash111/ecom/internal/models"
	"github.com/abhilash111/ecom/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService services.AuthService
	otpService  services.OTPService
}

func NewAuthController(authService services.AuthService, otpService services.OTPService) *AuthController {
	return &AuthController{
		authService: authService,
		otpService:  otpService,
	}
}

func (c *AuthController) Register(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Print User Object for debugging
	fmt.Println("Registering User:", user.Password)

	if err := c.authService.Register(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (c *AuthController) LoginWithEmail(ctx *gin.Context) {
	var credentials struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.authService.LoginWithEmail(credentials.Email, credentials.Password)
	fmt.Println("User found:", user)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := c.authService.GenerateToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (c *AuthController) RequestOTP(ctx *gin.Context) {
	var request struct {
		PhoneNumber string `json:"phone_number" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	otp, err := c.otpService.GenerateOTP(request.PhoneNumber)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate OTP"})
		return
	}

	// In production, you would send the OTP via SMS
	ctx.JSON(http.StatusOK, gin.H{"message": "OTP generated", "otp": otp.Code})
}

func (c *AuthController) LoginWithOTP(ctx *gin.Context) {
	var request struct {
		PhoneNumber string `json:"phone_number" binding:"required"`
		OTP         string `json:"otp" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	valid, err := c.otpService.VerifyOTP(request.PhoneNumber, request.OTP)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify OTP"})
		return
	}

	if !valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
		return
	}

	user, err := c.authService.GetUserByPhone(request.PhoneNumber)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	token, err := c.authService.GenerateToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "OTP verified successfully",
		"token":   token,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

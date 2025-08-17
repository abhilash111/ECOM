package controllers

import (
	"net/http"
	"time"

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
	var request models.RegisterRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.authService.RegisterUser(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenPair, err := c.authService.GenerateTokenPair(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate tokens"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":       "User registered successfully",
		"access_token":  tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
	})
}

func (c *AuthController) LoginWithEmail(ctx *gin.Context) {
	var request models.EmailLoginRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.authService.LoginWithEmail(request.Email, request.Password)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	tokenPair, err := c.authService.GenerateTokenPair(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate tokens"})
		return
	}

	session, err := c.authService.CreateSession(user, ctx.GetHeader("User-Agent"), ctx.ClientIP())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create session"})
		return
	}

	ctx.SetCookie("session_id", session.SessionID, int(time.Until(session.ExpiresAt).Seconds()), "/", "", false, true)
	userResponse := models.UserResponse{
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		Role:        string(user.Role),
	}
	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
		"user":          userResponse,
	})
}

func (c *AuthController) RequestOTP(ctx *gin.Context) {
	var request models.OTPRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	otp, err := c.otpService.GenerateOTP(request.PhoneNumber)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "OTP generated successfully",
		"otp":     otp.Code, // Remove in production
	})
}

func (c *AuthController) LoginWithOTP(ctx *gin.Context) {
	var request models.VerifyOTPRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	valid, err := c.otpService.VerifyOTP(request.PhoneNumber, request.OTP)
	if err != nil || !valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid OTP"})
		return
	}

	user, err := c.authService.FindUserByPhone(request.PhoneNumber)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	tokenPair, err := c.authService.GenerateTokenPair(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate tokens"})
		return
	}

	session, err := c.authService.CreateSession(user, ctx.GetHeader("User-Agent"), ctx.ClientIP())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create session"})
		return
	}

	ctx.SetCookie("session_id", session.SessionID, int(time.Until(session.ExpiresAt).Seconds()), "/", "", false, true)

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
	})
}

func (c *AuthController) Logout(ctx *gin.Context) {
	var request struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.authService.RevokeRefreshToken(request.RefreshToken); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to revoke token"})
		return
	}

	sessionID, err := ctx.Cookie("session_id")
	if err == nil {
		_ = c.authService.DeleteSession(sessionID)
		ctx.SetCookie("session_id", "", -1, "/", "", false, true)
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

func (c *AuthController) RefreshToken(ctx *gin.Context) {
	var request struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenPair, err := c.authService.RefreshAccessToken(request.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, tokenPair)
}

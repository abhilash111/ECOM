package user

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/abhilash111/ecom/config"
	"github.com/abhilash111/ecom/internal/auth"
	"github.com/abhilash111/ecom/internal/types"
	"github.com/abhilash111/ecom/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/login", h.handleLogin)
	router.POST("/register", h.handleRegister)

	authGroup := router.Group("/users")
	authGroup.Use(auth.JWTAuthMiddleware(h.store)) // Apply JWT auth middleware to this group
	// admin routes
	authGroup.GET("/:userID", h.handleGetUser)
}

func (h *Handler) handleLogin(c *gin.Context) {
	var user types.LoginUserPayload
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := utils.Validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid payload: %v", errors)})
		return
	}

	u, err := h.store.GetUserByEmail(user.Email)
	if err != nil || !auth.ComparePasswords(u.Password, []byte(user.Password)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email or password"})
		return
	}

	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, u.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *Handler) handleRegister(c *gin.Context) {
	var user types.RegisterUserPayload
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := utils.Validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid payload: %v", errors)})
		return
	}

	_, err := h.store.GetUserByEmail(user.Email)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user with this email already exists"})
		return
	}

	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	err = h.store.CreateUser(types.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  hashedPassword,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}

func (h *Handler) handleGetUser(c *gin.Context) {
	userIDStr := c.Param("userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := h.store.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

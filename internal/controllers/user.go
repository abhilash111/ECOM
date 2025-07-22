package controllers

import (
	"net/http"
	"strings"

	"github.com/abhilash111/ecom/internal/models"
	"github.com/abhilash111/ecom/internal/services"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{userService: userService}
}

func (c *UserController) GetCurrentUser(ctx *gin.Context) {
	username, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user, err := c.userService.GetUserByEmail(username.(string))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Convert roles string to slice
	var roles []models.Role
	if user.Roles != "" {
		roleStrs := strings.Split(user.Roles, ",")
		for _, r := range roleStrs {
			roles = append(roles, models.Role(r))
		}
	}

	response := models.UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Roles:       roles,
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *UserController) GetAllUsers(ctx *gin.Context) {
	// Implementation for admin-only endpoint
	ctx.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"})
}

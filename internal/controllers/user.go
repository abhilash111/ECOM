package controllers

import (
	"fmt"
	"net/http"

	"github.com/abhilash111/ecom/internal/services"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{userService: userService}
}

func (c *UserController) GetProfile(ctx *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := ctx.Get("userID")
	fmt.Println("User ID from context:", userID, "Exists:", exists)
	if !exists {
		fmt.Println("User ID not found in context")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// // Convert userID to uint
	// id, ok := userID.(float64)
	// if !ok {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
	// 	return
	// }

	user, err := c.userService.GetUserByID(userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":          user.ID,
			"name":        user.Name,
			"email":       user.Email,
			"phoneNumber": user.PhoneNumber,
			"role":        user.Role,
		},
	})
}

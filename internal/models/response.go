package models

type UserResponse struct {
	Name        string `json:"name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Role        string `json:"role" binding:"required"`
}

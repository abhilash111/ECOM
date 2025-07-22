package models

import "gorm.io/gorm"

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type User struct {
	gorm.Model
	Email       string `gorm:"unique"`
	PhoneNumber string `gorm:"unique"`
	CognitoID   string `gorm:"unique"`
	Roles       string // Comma-separated roles (e.g., "admin,user")
}

type UserResponse struct {
	ID          uint   `json:"id"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Roles       []Role `json:"roles"`
}

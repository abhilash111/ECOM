package models

import "gorm.io/gorm"

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type User struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Email       string `gorm:"unique;not null"`
	PhoneNumber string `gorm:"unique;not null"`
	Password    string `gorm:"not null" json:"password"`
	Role        Role   `gorm:"type:user_role;default:'user'"`
	IsVerified  bool   `gorm:"default:false"`
}

package models

import (
	"time"

	"gorm.io/gorm"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshToken struct {
	gorm.Model
	UserID    uint   `gorm:"not null;index"`
	Token     string `gorm:"unique;not null"`
	ExpiresAt time.Time
	Revoked   bool
	UserAgent string
	IPAddress string
}

type Session struct {
	gorm.Model
	UserID    uint   `gorm:"not null;index"`
	SessionID string `gorm:"unique;not null"`
	ExpiresAt time.Time
	UserAgent string
	IPAddress string
}

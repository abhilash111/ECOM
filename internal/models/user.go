package models

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type User struct {
	gorm.Model
	Name          string `gorm:"not null"`
	PhoneNumber   string `gorm:"unique;not null"`
	Email         string `gorm:"unique"`
	Password      string
	IsPassword    bool               `gorm:"default:false"`
	Role          Role               `gorm:"type:role;default:'user'"`
	Subscriptions []UserSubscription `gorm:"foreignKey:UserID"`
}

type PackageType string

const (
	ConsumerPremium PackageType = "consumer_premium"
	ConsumerBasic   PackageType = "consumer_basic"
	AgentPremium    PackageType = "agent_premium"
	AgentBasic      PackageType = "agent_basic"
	FreeTrial       PackageType = "free_trial"
)

type Package struct {
	gorm.Model
	Name             PackageType `gorm:"type:package_type;unique;not null"`
	ViewListingLimit *int        // NULL means unlimited
	AddListingLimit  int
	Price            float64
}

type UserSubscription struct {
	gorm.Model
	UserID               uint      `gorm:"not null"`
	PackageID            uint      `gorm:"not null"`
	StartsAt             time.Time `gorm:"not null"`
	ExpiresAt            time.Time `gorm:"not null"`
	RemainingViewCredits *int
	RemainingAddCredits  int
	User                 User    `gorm:"foreignKey:UserID"`
	Package              Package `gorm:"foreignKey:PackageID"`
}

package models

type RegisterRequest struct {
	Name             string      `json:"name" binding:"required"`
	PhoneNumber      string      `json:"phone_number" binding:"required"`
	Email            string      `json:"email" binding:"required,email"`
	Password         string      `json:"password"`
	SubscriptionPack PackageType `json:"subscription_pack" binding:"required"`
}

type EmailLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type OTPRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type VerifyOTPRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	OTP         string `json:"otp" binding:"required"`
}

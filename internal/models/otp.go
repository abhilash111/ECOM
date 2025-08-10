package models

type OTP struct {
	PhoneNumber string `json:"phone_number"`
	Code        string `json:"code"`
	ExpiresAt   int64  `json:"expires_at"`
}

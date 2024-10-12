package models

type User struct {
	ID         int64  `json:"id"`
	Email      string `json:"email"`
	Password   string `json:"-"` // Store hashed password
	OTP        string `json:"otp,omitempty"`
	IsVerified bool   `json:"is_verified"`
}

package models

type User struct {
	ID                  int64  `json:"id"`
	Company_Name        string `json:"company_name"`
	Company_Category_Id *int64 `json:"category_id"`
	Email               string `json:"email"`
	Password            string `json:"password"`
	OTP                 string `json:"otp,omitempty"`
	IsVerified          bool   `json:"is_verified"`
	Token               string `json:"token"`
}

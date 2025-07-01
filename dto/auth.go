package dto

type AuthLogin struct {
	Email    string `form:"email" json:"email" db:"email" binding:"required,email"`
	Password string `form:"password" json:"password" db:"password" binding:"required"`
}

type AuthRegister struct {
	Name            string `form:"name" json:"name" db:"name" binding:"required"`
	Email           string `form:"email" json:"email" db:"email" binding:"required,email"`
	PhoneNumber     string `form:"phoneNumber" json:"phoneNumber" db:"phone_number" binding:"required"`
	Password        string `form:"password" json:"password" db:"password" binding:"required"`
	ConfirmPassword string `form:"confPass" json:"confPass" binding:"required"`
}

type AuthResetPass struct {
	Email    string `json:"email" form:"email"`
	OTP      string `json:"otp" form:"otp"`
	NewPass  string `json:"newPass" form:"newPass"`
	ConfPass string `json:"confPass" form:"confPass"`
}

package utils

type Response struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	Errors   any    `json:"errors,omitempty"`
	PageInfo any    `json:"pageInfo,omitempty"`
	Result   any    `json:"results,omitempty"`
}

type ResponseUser struct {
	Name        string `form:"name" json:"name" db:"name" binding:"required"`
	Email       string `form:"email" json:"email" db:"email" binding:"required,email"`
	PhoneNumber string `form:"phoneNumber" json:"phoneNumber" db:"phone_number" binding:"required"`
}

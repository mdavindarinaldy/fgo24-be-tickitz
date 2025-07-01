package utils

type Response struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	Errors   any    `json:"errors,omitempty"`
	PageInfo any    `json:"pageInfo,omitempty"`
	Result   any    `json:"results,omitempty"`
}

type ResponseUser struct {
	Name        string `form:"name,omitempty" json:"name,omitempty"`
	Email       string `form:"email,omitempty" json:"email,omitempty"`
	PhoneNumber string `form:"phoneNumber,omitempty" json:"phoneNumber,omitempty"`
}

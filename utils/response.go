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

type PageData struct {
	TotalData   int `json:"totalData"`
	TotalPage   int `json:"totalPage"`
	CurrentPage int `json:"currentPage"`
}

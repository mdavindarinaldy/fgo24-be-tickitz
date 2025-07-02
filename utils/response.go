package utils

import "time"

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

type ResponseMovie struct {
	Title       string    `json:"title" db:"title"`
	Synopsis    string    `json:"synopsis" db:"synopsis"`
	ReleaseDate time.Time `json:"releaseDate" db:"release_date"`
	Price       float64   `json:"price" db:"price"`
	Runtime     int       `json:"runtime" db:"runtime"`
	Poster      string    `json:"poster" db:"poster"`
	Backdrop    string    `json:"backdrop" db:"backdrop"`
}

type PageData struct {
	TotalData   int `json:"totalData"`
	TotalPage   int `json:"totalPage"`
	CurrentPage int `json:"currentPage"`
}

package models

import "time"

type User struct {
	Id          int       `db:"id" json:"id"`
	Name        string    `form:"name" json:"name" db:"name" binding:"required"`
	Email       string    `form:"email" json:"email" db:"email" binding:"required,email"`
	PhoneNumber string    `form:"phoneNumber" json:"phoneNumber" db:"phone_number" binding:"required"`
	Password    string    `form:"password" json:"password" db:"password" binding:"required"`
	Role        string    `form:"role" json:"role" db:"role" binding:"required"`
	CreatedAt   time.Time `form:"createdAt" json:"createdAt" db:"created_at" binding:"required"`
	UpdatedAt   time.Time `form:"updatedAt" json:"updatedAt" db:"updated_at" binding:"required"`
}

type ResponseUser struct {
	Id          int    `db:"id" json:"id"`
	Name        string `form:"name" json:"name" db:"name" binding:"required"`
	Email       string `form:"email" json:"email" db:"email" binding:"required,email"`
	PhoneNumber string `form:"phoneNumber" json:"phoneNumber" db:"phone_number" binding:"required"`
}

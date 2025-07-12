package dto

type UpdateUserRequest struct {
	Email           *string `form:"email" json:"email" binding:"omitempty,email"`
	Password        *string `form:"password" json:"password" binding:"omitempty,min=6"`
	ConfirmPassword *string `form:"confirmPassword" json:"confirmPassword" binding:"omitempty"`
	Name            *string `form:"fullname" json:"fullname" binding:"omitempty"`
	PhoneNumber     *string `form:"phoneNumber" json:"phoneNumber" binding:"omitempty"`
	ProfilePicture  *string `form:"profilePicture" json:"profilePicture" binding:"omitempty"`
}

type UpdateUserResult struct {
	Email          string  `json:"email" db:"email"`
	Name           string  `json:"name" db:"name"`
	PhoneNumber    string  `json:"phoneNumber" db:"phone_number"`
	ProfilePicture *string `json:"profilePicture" db:"profile_picture"`
}

type Profile struct {
	Email          string  `json:"email" db:"email"`
	Name           string  `json:"name" db:"name"`
	Role           string  `json:"role" db:"role"`
	PhoneNumber    string  `json:"phoneNumber" db:"phone_number"`
	ProfilePicture *string `json:"profilePicture" db:"profile_picture"`
}

type CheckPass struct {
	Password string `json:"password" db:"password" binding:"required"`
}

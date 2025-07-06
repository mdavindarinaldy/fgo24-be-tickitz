package dto

type UpdateUserRequest struct {
	Email           *string `json:"email" binding:"omitempty,email"`
	Password        *string `json:"password" binding:"omitempty,min=6"`
	ConfirmPassword *string `json:"confirmPassword" binding:"eqfield=Password"`
	Name            *string `json:"name" binding:"omitempty"`
	PhoneNumber     *string `json:"phoneNumber" binding:"omitempty"`
	ProfilePicture  *string `json:"profilePicture" binding:"omitempty,uri"`
}

type UpdateUserResult struct {
	Email          string  `json:"email" db:"email"`
	Name           string  `json:"name" db:"name"`
	PhoneNumber    string  `json:"phoneNumber" db:"phone_number"`
	ProfilePicture *string `json:"profilePicture" db:"profile_picture"`
}

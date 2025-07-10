package models

import (
	"be-tickitz/dto"
	"be-tickitz/utils"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/matthewhartstonge/argon2"
)

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

type UserCredentials struct {
	Id        int       `db:"id" json:"id"`
	Email     string    `form:"email" json:"email" db:"email" binding:"required,email"`
	Password  string    `form:"password" json:"password" db:"password" binding:"required"`
	Role      string    `form:"role" json:"role" db:"role" binding:"required"`
	CreatedAt time.Time `form:"createdAt" json:"createdAt" db:"created_at" binding:"required"`
	UpdatedAt time.Time `form:"updatedAt" json:"updatedAt" db:"updated_at" binding:"required"`
}

func UpdateUserData(userId int, request dto.UpdateUserRequest) (dto.UpdateUserResult, error) {
	conn, err := utils.DBConnect()
	if err != nil {
		return dto.UpdateUserResult{}, err
	}
	defer conn.Close()

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return dto.UpdateUserResult{}, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		} else {
			tx.Commit(context.Background())
		}
	}()

	if request.Password != nil && *request.Password != *request.ConfirmPassword {
		return dto.UpdateUserResult{}, errors.New("password and confirm password do not match")
	}

	if request.Email != nil && request.Password != nil {
		argon := argon2.DefaultConfig()
		hashedPassword, err := argon.HashEncoded([]byte(*request.Password))
		if err != nil {
			return dto.UpdateUserResult{}, err
		}
		_, err = tx.Exec(context.Background(),
			`UPDATE users SET email = $1, password = $2 WHERE id = $3`,
			*request.Email, string(hashedPassword), userId)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				return dto.UpdateUserResult{}, errors.New("email already used by another user")
			}
			return dto.UpdateUserResult{}, err
		}
	} else if request.Email != nil {
		_, err = tx.Exec(context.Background(),
			`UPDATE users SET email = $1 WHERE id = $2`,
			*request.Email, userId)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				return dto.UpdateUserResult{}, errors.New("email already used by another user")
			}
			return dto.UpdateUserResult{}, err
		}
	} else if request.Password != nil {
		argon := argon2.DefaultConfig()
		hashedPassword, err := argon.HashEncoded([]byte(*request.Password))
		if err != nil {
			return dto.UpdateUserResult{}, err
		}
		_, err = tx.Exec(context.Background(),
			`UPDATE users SET password = $1 WHERE id = $2`,
			string(hashedPassword), userId)
		if err != nil {
			return dto.UpdateUserResult{}, err
		}
	}

	if request.Name != nil && request.PhoneNumber != nil && request.ProfilePicture != nil {
		_, err = tx.Exec(context.Background(),
			`UPDATE profiles SET name = $1, phone_number = $2, profile_picture = $3 WHERE id_user = $4`,
			*request.Name, *request.PhoneNumber, *request.ProfilePicture, userId)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				return dto.UpdateUserResult{}, errors.New("phone number already used by another user")
			}
			return dto.UpdateUserResult{}, err
		}
	} else if request.Name != nil && request.PhoneNumber != nil {
		_, err = tx.Exec(context.Background(),
			`UPDATE profiles SET name = $1, phone_number = $2 WHERE id_user = $3`,
			*request.Name, *request.PhoneNumber, userId)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				return dto.UpdateUserResult{}, errors.New("phone number already used by another user")
			}
			return dto.UpdateUserResult{}, err
		}
	} else if request.Name != nil && request.ProfilePicture != nil {
		_, err = tx.Exec(context.Background(),
			`UPDATE profiles SET name = $1, profile_picture = $2 WHERE id_user = $3`,
			*request.Name, *request.ProfilePicture, userId)
		if err != nil {
			return dto.UpdateUserResult{}, err
		}
	} else if request.PhoneNumber != nil && request.ProfilePicture != nil {
		_, err = tx.Exec(context.Background(),
			`UPDATE profiles SET phone_number = $1, profile_picture = $2 WHERE id_user = $3`,
			*request.PhoneNumber, *request.ProfilePicture, userId)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				return dto.UpdateUserResult{}, errors.New("phone number already used by another user")
			}
			return dto.UpdateUserResult{}, err
		}
	} else if request.Name != nil {
		_, err = tx.Exec(context.Background(),
			`UPDATE profiles SET name = $1 WHERE id_user = $2`,
			*request.Name, userId)
		if err != nil {
			return dto.UpdateUserResult{}, err
		}
	} else if request.PhoneNumber != nil {
		_, err = tx.Exec(context.Background(),
			`UPDATE profiles SET phone_number = $1 WHERE id_user = $2`,
			*request.PhoneNumber, userId)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				return dto.UpdateUserResult{}, errors.New("phone number already used by another user")
			}
			return dto.UpdateUserResult{}, err
		}
	} else if request.ProfilePicture != nil {
		_, err = tx.Exec(context.Background(),
			`UPDATE profiles SET profile_picture = $1 WHERE id_user = $2`,
			*request.ProfilePicture, userId)
		if err != nil {
			return dto.UpdateUserResult{}, err
		}
	}

	var userData dto.UpdateUserResult

	err = tx.QueryRow(context.Background(),
		`SELECT u.email, p.name, p.phone_number, p.profile_picture
         FROM users u
         JOIN profiles p ON u.id = p.id_user
         WHERE u.id = $1`, userId).
		Scan(&userData.Email, &userData.Name, &userData.PhoneNumber, &userData.ProfilePicture)

	if err != nil {
		return dto.UpdateUserResult{}, err
	}

	return userData, nil
}

func GetProfileUser(userId int) (dto.Profile, error) {
	conn, err := utils.DBConnect()
	if err != nil {
		return dto.Profile{}, err
	}
	defer conn.Close()
	row, err := conn.Query(context.Background(), `
		SELECT p.name, u.email, u.role, p.phone_number, p.profile_picture 
		FROM profiles p
		JOIN users u ON u.id = p.id_user
		WHERE p.id=$1`, userId)
	if err != nil {
		return dto.Profile{}, err
	}
	result, err := pgx.CollectOneRow[dto.Profile](row, pgx.RowToStructByName)
	if err != nil {
		return dto.Profile{}, err
	}
	return result, nil
}

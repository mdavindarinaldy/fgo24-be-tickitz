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

	var hashedPassword *string
	if request.Password != nil {
		argon := argon2.DefaultConfig()
		hash, err := argon.HashEncoded([]byte(*request.Password))
		if err != nil {
			return dto.UpdateUserResult{}, err
		}
		hashStr := string(hash)
		hashedPassword = &hashStr
	}

	_, err = tx.Exec(context.Background(),
		`UPDATE users
		 SET email = COALESCE($1, email),
		     password = COALESCE($2, password),
			 updated_at = $3
		 WHERE id = $4`,
		request.Email, hashedPassword, time.Now(), userId)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return dto.UpdateUserResult{}, errors.New("email already used by another user")
		}
		return dto.UpdateUserResult{}, err
	}

	_, err = tx.Exec(context.Background(),
		`UPDATE profiles
		 SET name = COALESCE($1, name),
		     phone_number = COALESCE($2, phone_number),
		     profile_picture = COALESCE($3, profile_picture),
			 updated_at = $4
		 WHERE id_user = $5`,
		request.Name, request.PhoneNumber, request.ProfilePicture, time.Now(), userId)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return dto.UpdateUserResult{}, errors.New("phone number already used by another user")
		}
		return dto.UpdateUserResult{}, err
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
		WHERE p.id_user=$1`, userId)
	if err != nil {
		return dto.Profile{}, err
	}
	result, err := pgx.CollectOneRow[dto.Profile](row, pgx.RowToStructByName)
	if err != nil {
		return dto.Profile{}, err
	}
	return result, nil
}

func CheckPass(userId int, pass string) (bool, error) {
	conn, err := utils.DBConnect()
	if err != nil {
		return false, err
	}
	defer conn.Close()
	row, err := conn.Query(context.Background(), `
		SELECT password 
		FROM users
		WHERE id=$1`, userId)
	if err != nil {
		return false, err
	}
	result, err := pgx.CollectOneRow[dto.CheckPass](row, pgx.RowToStructByName)
	if err != nil {
		return false, err
	}
	ok, _ := argon2.VerifyEncoded([]byte(pass), []byte(result.Password))
	return ok, nil
}

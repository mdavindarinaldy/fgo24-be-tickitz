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

func HandleRegister(user dto.AuthRegister) error {
	if user.Email == "" || user.Name == "" || user.Password == "" || user.PhoneNumber == "" {
		return errors.New("user data should not be empty")
	}
	if user.Password != user.ConfirmPassword {
		return errors.New("password and confirm password doesn't match")
	}

	argon := argon2.DefaultConfig()
	hashedPassword, err := argon.HashEncoded([]byte(user.Password))
	if err != nil {
		return errors.New("failed to hash password")
	}

	conn, err := utils.DBConnect()
	if err != nil {
		return err
	}
	defer conn.Close()

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		} else {
			tx.Commit(context.Background())
		}
	}()

	var userId int
	err = tx.QueryRow(context.Background(),
		`INSERT INTO users (email, password, role, created_at)
         VALUES ($1, $2, 'user', $3)
         RETURNING id`,
		user.Email, string(hashedPassword), time.Now()).Scan(&userId)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate key") {
			return errors.New("email already used by another user")
		}
		return err
	}

	_, err = tx.Exec(context.Background(),
		`INSERT INTO profiles (id_user, name, phone_number, created_at)
         VALUES ($1, $2, $3, $4)`,
		userId, user.Name, user.PhoneNumber, time.Now())
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate key") {
			return errors.New("phone number already used by another user")
		}
		return err
	}

	return nil
}

func GetUser(email string) (UserCredentials, error) {
	conn, err := utils.DBConnect()
	if err != nil {
		return UserCredentials{}, err
	}
	defer conn.Close()

	rows, err := conn.Query(context.Background(),
		`SELECT * FROM users WHERE email=$1`, email)
	if err != nil {
		return UserCredentials{}, err
	}

	userData, err := pgx.CollectOneRow[UserCredentials](rows, pgx.RowToStructByName)
	if err != nil {
		return UserCredentials{}, err
	}

	return userData, nil
}

func ResetPass(id int, newPass string) error {
	argon := argon2.DefaultConfig()
	hashedPassword, err := argon.HashEncoded([]byte(newPass))
	if err != nil {
		return err
	}

	conn, err := utils.DBConnect()
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Exec(
		context.Background(),
		`UPDATE users SET password = $1 
		WHERE id = $2`, hashedPassword, id)
	if err != nil {
		return err
	}
	return nil
}

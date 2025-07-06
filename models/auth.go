package models

import (
	"be-tickitz/dto"
	"be-tickitz/utils"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

func HandleRegister(user dto.AuthRegister) error {
	if user.Email == "" || user.Name == "" || user.Password == "" || user.PhoneNumber == "" {
		return errors.New("user data should not be empty")
	}
	if user.Password != user.ConfirmPassword {
		return errors.New("password and confirm password doesn't match")
	}
	conn, err := utils.DBConnect()
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Exec(
		context.Background(),
		`
		INSERT INTO users (name, email, phone_number, password, created_at, role)
		VALUES
		($1,$2,$3,$4,$5,'user');
		`,
		user.Name, user.Email, user.PhoneNumber, user.Password, time.Now())
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate key") {
			return errors.New("email already used by another user")
		}
		return err
	}
	return nil
}

func GetUser(email string) (User, error) {
	conn, err := utils.DBConnect()
	if err != nil {
		return User{}, err
	}
	defer conn.Close()

	rows, err := conn.Query(context.Background(),
		`SELECT * FROM users WHERE email=$1`, email)
	if err != nil {
		return User{}, err
	}

	userData, err := pgx.CollectOneRow[User](rows, pgx.RowToStructByName)
	if err != nil {
		return User{}, err
	}

	return userData, nil
}

func ResetPass(id int, newPass string) error {
	conn, err := utils.DBConnect()
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Exec(
		context.Background(),
		`UPDATE users SET password = $1 
		WHERE id = $2`, id, newPass)
	if err != nil {
		return err
	}
	return nil
}

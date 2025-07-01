package models

import (
	"be-tickitz/dto"
	"be-tickitz/utils"
	"context"
	"errors"
	"strings"
	"time"
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
		($1,$2,$3,$4,$5,"user");
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

package controllers

import (
	"be-tickitz/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(user models.UserCredentials) (string, error) {
	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.Id,
		"role":   user.Role,
		"iat":    time.Now().Unix(),
		"exp":    time.Now().Add(1 * time.Hour).Unix(),
	})
	token, err := generateToken.SignedString([]byte(os.Getenv("APP_SECRET")))
	if err != nil {
		return token, err
	}
	return token, nil
}

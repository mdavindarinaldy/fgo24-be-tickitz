package controllers

import (
	"be-tickitz/dto"
	"be-tickitz/models"
	"be-tickitz/utils"
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/matthewhartstonge/argon2"
)

// AuthRegister handles user registration
// @Summary Register a new user
// @Description Register a new user with name, email, phone number, and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body dto.AuthRegister true "User registration data"
// @Success 201 {object} utils.Response{result=utils.ResponseUser} "User registered successfully"
// @Failure 400 {object} utils.Response{errors=string} "Bad request (e.g., email already used, empty data, or password mismatch)"
// @Failure 500 {object} utils.Response{errors=string} "Internal server error"
// @Router /auth/register [post]
func AuthRegister(c *gin.Context) {
	user := dto.AuthRegister{}
	c.ShouldBind(&user)
	err := models.HandleRegister(user)
	if err != nil {
		if err.Error() == "email already used by another user" || err.Error() == "phone number already used by another user" || err.Error() == "user data should not be empty" || err.Error() == "password and confirm password doesn't match" {
			c.JSON(http.StatusBadRequest, utils.Response{
				Success: false,
				Message: "Failed to register user",
				Errors:  err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Internal server error",
			Errors:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, utils.Response{
		Success: true,
		Message: "Create user success!",
		Result: utils.ResponseUser{
			Name:        user.Name,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
		},
	})
}

// AuthLogin handles user login
// @Summary User login
// @Description Authenticate a user with email and password, returning a token
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body dto.AuthLogin true "User login credentials"
// @Success 200 {object} utils.Response{result=string} "Login successful with token"
// @Failure 400 {object} utils.Response{errors=string} "Bad request (e.g., user not registered, wrong password, or token generation failed)"
// @Failure 500 {object} utils.Response{errors=string} "Internal server error"
// @Router /auth/login [post]
func AuthLogin(c *gin.Context) {
	user := dto.AuthLogin{}
	c.ShouldBind(&user)
	userData, err := models.GetUser(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Internal server error",
			Errors:  err.Error(),
		})
		return
	}
	if userData == (models.UserCredentials{}) {
		c.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "User is not registered",
		})
		return
	}

	ok, err := argon2.VerifyEncoded([]byte(user.Password), []byte(userData.Password))
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Internal server error",
			Errors:  err.Error(),
		})
		return
	}
	if !ok {
		c.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Wrong password",
		})
		return
	}

	token, err := GenerateToken(userData)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Login success!",
		Result:  token,
	})
}

type email struct {
	Email string `json:"email" form:"email"`
}

// AuthForgotPass handles password reset request
// @Summary Request password reset
// @Description Send an OTP to the user's email for password reset
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body email true "User email"
// @Success 200 {object} utils.Response "OTP sent successfully"
// @Failure 400 {object} utils.Response "Bad request (e.g., email not registered)"
// @Failure 500 {object} utils.Response{errors=string} "Internal server error"
// @Router /auth/pass [post]
func AuthForgotPass(c *gin.Context) {
	emailUser := email{}
	c.ShouldBind(&emailUser)
	userData, err := models.GetUser(emailUser.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Internal server error",
			Errors:  err.Error(),
		})
		return
	}
	if userData != (models.UserCredentials{}) {
		rdClient := utils.RedisConnect()
		otp := fmt.Sprint(rand.Intn(900) + 100)
		endpoint := fmt.Sprintf("/auth/otp/%d", userData.Id)
		rdClient.Set(context.Background(), endpoint, otp, (3 * time.Minute))
		c.JSON(http.StatusOK, utils.Response{
			Success: true,
			Message: "OTP has been sent to your email that will expires within 3 minutes",
		})
		fmt.Printf("[EMAIL SIMULATION]\n OTP Code for %s to reset password is %s\n", emailUser.Email, otp)
		rdClient.Close()
	} else {
		c.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Email is not registered!",
		})
	}
}

// AuthResetPass handles password reset
// @Summary Reset user password
// @Description Reset user password using email, OTP, and new password
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body dto.AuthResetPass true "Password reset data"
// @Success 200 {object} utils.Response{result=utils.ResponseUser} "Password reset successful"
// @Failure 400 {object} utils.Response "Bad request (e.g., invalid OTP, password mismatch, or email not registered)"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /auth/pass [patch]
func AuthResetPass(c *gin.Context) {
	credentials := dto.AuthResetPass{}
	c.ShouldBind(&credentials)
	userData, err := models.GetUser(credentials.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Internal server error",
			Errors:  err.Error(),
		})
		return
	}
	if userData != (models.UserCredentials{}) {
		rdClient := utils.RedisConnect()
		redisEndpoint := fmt.Sprintf("/auth/otp/%d", userData.Id)
		data := rdClient.Get(context.Background(), redisEndpoint)
		if data.Val() == credentials.OTP {
			if credentials.NewPass == credentials.ConfPass {
				err := models.ResetPass(userData.Id, credentials.NewPass)
				if err != nil {
					c.JSON(http.StatusInternalServerError, utils.Response{
						Success: false,
						Message: "Internal server error",
						Errors:  err.Error(),
					})
				} else {
					c.JSON(http.StatusOK, utils.Response{
						Success: true,
						Message: "Change pass success!",
						Result: utils.ResponseUser{
							Email: userData.Email,
						},
					})
				}
			} else {
				c.JSON(http.StatusBadRequest, utils.Response{
					Success: false,
					Message: "Password and confirm password does not match",
				})
			}
		} else {
			c.JSON(http.StatusBadRequest, utils.Response{
				Success: false,
				Message: "OTP is wrong or already expired!",
			})
		}
		defer rdClient.Close()
	}
}

// Logout adds the user's token to the blacklist in Redis with expiration based on token's exp claim
// @Summary Logout user
// @Description Invalidates the user's JWT by adding it to a Redis blacklist with expiration based on the token's exp claim
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response "Successful logout"
// @Failure 400 {object} utils.Response{errors=string} "Bad request due to missing or invalid token"
// @Failure 500 {object} utils.Response{errors=string} "Internal server error"
// @Router /auth/logout [post]
func AuthLogout(c *gin.Context) {
	secretKey := os.Getenv("APP_SECRET")
	token := strings.Split(c.GetHeader("Authorization"), "Bearer ")
	rawToken, _ := jwt.Parse(token[1], func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	var expiresAt time.Time = time.Unix(int64(rawToken.Claims.(jwt.MapClaims)["exp"].(float64)), 0)
	var duration time.Duration = time.Until(expiresAt)
	if duration <= 0 {
		c.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Token already expired",
		})
		return
	}

	redisClient := utils.RedisConnect()
	defer redisClient.Close()

	err := redisClient.SetEx(context.Background(), "blacklist:"+token[1], "true", duration).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to logout",
			Errors:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Successfully logged out",
	})
}

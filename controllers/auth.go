package controllers

import (
	"be-tickitz/dto"
	"be-tickitz/models"
	"be-tickitz/utils"
	"context"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthRegister(c *gin.Context) {
	user := dto.AuthRegister{}
	c.ShouldBind(&user)
	err := models.HandleRegister(user)
	if err != nil {
		if err.Error() == "email already used by another user" || err.Error() == "user data should not be empty" || err.Error() == "password and confirm password doesn't match" {
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
	if userData == (models.User{}) {
		c.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "User is not registered",
		})
		return
	}
	if userData.Password != user.Password {
		c.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Wrong password",
		})
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Login success!",
		Result: utils.ResponseUser{
			Name:        userData.Name,
			Email:       userData.Email,
			PhoneNumber: userData.PhoneNumber,
		},
	})
}

func AuthForgotPass(c *gin.Context) {
	type email struct {
		Email string `json:"email" form:"email"`
	}
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
	if userData != (models.User{}) {
		rdClient := utils.RedisConnect()
		otp := fmt.Sprint(rand.Intn(900) + 100)
		endpoint := fmt.Sprintf("/auth/otp/%d", userData.Id)
		rdClient.Set(context.Background(), endpoint, otp, (3 * 60))
		c.JSON(http.StatusOK, utils.Response{
			Success: true,
			Message: "OTP has been sent to your email that will expires within 3 minutes",
		})
		fmt.Printf("[EMAIL SIMULATION]\n OTP Code for %s to reset password is %s\n", userData.Email, otp)
		rdClient.Close()
	} else {
		c.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Email is not registered!",
		})
	}
}

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
	if userData != (models.User{}) {
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
					})
				} else {
					c.JSON(http.StatusOK, utils.Response{
						Success: true,
						Message: "Change pass success!",
						Result: utils.ResponseUser{
							Name:  userData.Name,
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
		rdClient.Close()
	}
}

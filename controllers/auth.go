package controllers

import (
	"be-tickitz/dto"
	"be-tickitz/models"
	"be-tickitz/utils"
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

}

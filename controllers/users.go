package controllers

import (
	"be-tickitz/dto"
	"be-tickitz/models"
	"be-tickitz/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UpdateUser updates user credentials and/or profile
// @Summary Update user data
// @Description Updates user credentials (email, password) and/or profile (name, phone number, profile picture)
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body dto.UpdateUserRequest true "User update data"
// @Success 200 {object} utils.Response{result=dto.UpdateUserResult} "User data updated successfully"
// @Failure 400 {object} utils.Response{errors=string} "Bad request (e.g., invalid input, email/phone already used, password mismatch)"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response{errors=string} "Internal server error"
// @Router /user [patch]
func UpdateUser(c *gin.Context) {
	userId, _ := c.Get("userId")
	role, _ := c.Get("role")
	if role != "user" {
		c.JSON(http.StatusUnauthorized, utils.Response{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}
	var request dto.UpdateUserRequest
	c.ShouldBindJSON(&request)

	updatedUser, err := models.UpdateUserProfile(int(userId.(float64)), request)
	if err != nil {
		if err.Error() == "email already used by another user" || err.Error() == "phone number already used by another user" || err.Error() == "password and confirm password do not match" {
			c.JSON(http.StatusBadRequest, utils.Response{
				Success: false,
				Message: "Failed to update user",
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

	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "User data updated successfully",
		Result:  updatedUser,
	})
}

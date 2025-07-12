package controllers

import (
	"be-tickitz/dto"
	"be-tickitz/models"
	"be-tickitz/utils"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UpdateUser updates user credentials and/or profile
// @Summary Update user data
// @Description Updates user credentials (email, password) and/or profile (name, phone number, profile picture)
// @Tags Profile
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param name formData string false "Name"
// @Param email formData string false "Email"
// @Param password formData string false "Password"
// @Param phone formData string false "Phone number"
// @Param file formData file false "Profile picture"
// @Success 200 {object} utils.Response{result=dto.UpdateUserResult} "User data updated successfully"
// @Failure 400 {object} utils.Response{errors=string} "Bad request (e.g., invalid input, email/phone already used, password mismatch)"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response{errors=string} "Internal server error"
// @Router /profile [patch]
func UpdateUser(c *gin.Context) {
	userId, _ := c.Get("userId")
	role, _ := c.Get("role")
	if role != "user" {
		c.JSON(http.StatusForbidden, utils.Response{
			Success: false,
			Message: "Forbidden",
		})
		return
	}
	var request dto.UpdateUserRequest
	c.ShouldBind(&request)

	file, _ := c.FormFile("file")
	fileName := ""
	if file != nil {
		if file.Size > 2*1024*1024 {
			c.JSON(http.StatusBadRequest, utils.Response{
				Success: false,
				Message: "File is too large",
			})
			return
		}
		ext := strings.ToLower(filepath.Ext(file.Filename))
		allowedExts := map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
		}
		if !allowedExts[ext] {
			c.JSON(http.StatusBadRequest, utils.Response{
				Success: false,
				Message: "Invalid file type. Only JPG, JPEG, PNG allowed",
			})
			return
		}
		fileExt := filepath.Ext(file.Filename)
		fileName = uuid.New().String() + fileExt
		err := c.SaveUploadedFile(file, "./uploads/profiles/"+fileName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, utils.Response{
				Success: false,
				Message: "Failed to save uploaded file",
			})
			return
		}
		request.ProfilePicture = &fileName
	}

	updatedUser, err := models.UpdateUserData(int(userId.(float64)), request)
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

// GetProfileUser retrieves user profile data
// @Summary Get profile user
// @Description Retrieve user profile data
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{result=[]dto.Profile} "Successful response with user profile data"
// @Failure 401 {object} utils.Response "Unauthorized access"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /profile [get]
func GetProfileUser(c *gin.Context) {
	userId, _ := c.Get("userId")
	user, err := models.GetProfileUser(int(userId.(float64)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Internal Server Error",
			Errors:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Success to get user's profile",
		Result:  user,
	})
}

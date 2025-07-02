package controllers

import (
	"be-tickitz/models"
	"be-tickitz/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetMovies(c *gin.Context) {
	search := strings.ToLower(c.DefaultQuery("search", ""))
	genre := strings.ToLower(c.DefaultQuery("genre", ""))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	movies, pageData, err := models.GetMovies(search, genre, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Success to get movies",
		PageInfo: utils.PageData{
			CurrentPage: pageData.CurrentPage,
			TotalPage:   pageData.TotalPage,
			TotalData:   pageData.TotalData,
		},
		Result: movies,
	})
}

func GetUpcomingMovies(c *gin.Context) {

}

package controllers

import (
	"be-tickitz/models"
	"be-tickitz/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetMovies retrieves a list of movies based on search and genre filters
// @Summary Get movies
// @Description Retrieve movies with optional search by title, genre filter, and pagination
// @Tags Movies
// @Accept json
// @Produce json
// @Param search query string false "Search by movie title (case-insensitive)" example:"inception"
// @Param genre query string false "Filter by genre (case-insensitive)" example:"action"
// @Param page query int false "Page number for pagination (default: 1)" example:"1"
// @Success 200 {object} utils.Response "Successful response with movies and pagination info"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /movies [get]
func GetMovies(c *gin.Context) {
	search := strings.ToLower(c.DefaultQuery("search", ""))
	genre := strings.ToLower(c.DefaultQuery("genre", ""))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	movies, pageData, err := models.GetMovies(search, genre, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Internal server error",
			Errors:  err.Error(),
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

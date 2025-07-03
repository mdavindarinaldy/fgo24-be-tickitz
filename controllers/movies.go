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

// GetDetailMovie retrieves details of a specific movie by ID
// @Summary Get movie details
// @Description Retrieve detailed information about a movie by its ID
// @Tags Movies
// @Accept json
// @Produce json
// @Param id path int true "Movie ID" example:"1"
// @Success 200 {object} utils.Response "Successful response with movie details"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /movies/{id} [get]
func GetDetailMovie(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	movie, err := models.GetMovie(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to get data",
		})
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Success to get data",
		Result:  movie,
	})
}

// @Summary Get upcoming movies
// @Description Retrieve a list of upcoming movies
// @Tags Movies
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response "Successful response with upcoming movies"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /movies/upcoming [get]
func GetUpcomingMovies(c *gin.Context) {
	movies, err := models.GetUpcomingMovies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to get data",
		})
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Success to get data",
		Result:  movies,
	})
}

// GetGenres retrieves the list of all genres
// @Summary Get all genres
// @Description Retrieve a list of all genres (admin only)
// @Tags Movies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{result=[]dto.SubData} "Successful response with genres list"
// @Failure 401 {object} utils.Response "Unauthorized access (requires admin role)"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /movies/genres [get]
func GetGenres(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusUnauthorized, utils.Response{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}
	genres, err := models.GetGenre()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Internal Server Error",
		})
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Success to get genres list",
		Result:  genres,
	})
}

// GetDirectors retrieves directors with optional search
// @Summary Get directors
// @Description Retrieve a list of directors with optional search by name (admin only)
// @Tags Movies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param search query string false "Search by director name (case-insensitive)" example:"nolan"
// @Success 200 {object} utils.Response{result=[]dto.SubData} "Successful response with directors list"
// @Failure 401 {object} utils.Response "Unauthorized access (requires admin role)"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /movies/directors [get]
func GetDirectors(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusUnauthorized, utils.Response{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}
	search := strings.ToLower(c.DefaultQuery("search", ""))
	directors, err := models.GetDirectors(search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Internal Server Error",
		})
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Success to get directors list",
		Result:  directors,
	})
}

// GetCasts retrieves casts with optional search
// @Summary Get casts
// @Description Retrieve a list of casts with optional search by name (admin only)
// @Tags Movies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param search query string false "Search by cast name (case-insensitive)" example:"dicaprio"
// @Success 200 {object} utils.Response{result=[]dto.SubData} "Successful response with casts list"
// @Failure 401 {object} utils.Response "Unauthorized access (requires admin role)"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /movies/casts [get]
func GetCasts(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusUnauthorized, utils.Response{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}
	search := strings.ToLower(c.DefaultQuery("search", ""))
	casts, err := models.GetCasts(search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Internal Server Error",
		})
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Success to get casts list",
		Result:  casts,
	})
}

func AddMovie(c *gin.Context) {
	// role, _ := c.Get("role")
	// if role != "admin" {
	// 	c.JSON(http.StatusUnauthorized, utils.Response{
	// 		Success: false,
	// 		Message: "Unauthorized",
	// 	})
	// 	return
	// }
	// userId, _ := c.Get("userId")
	// newMovie := dto.Movie{}
	// c.ShouldBind(&newMovie)
	// err := models.AddMovie(newMovie, int(userId.(float64)))
	// if err != nil {
	// 	if err.Error() == "new movie data should not be empty" {
	// 		c.JSON(http.StatusBadRequest, utils.Response{
	// 			Success: false,
	// 			Message: err.Error(),
	// 		})
	// 		return
	// 	}
	// 	c.JSON(http.StatusInternalServerError, utils.Response{
	// 		Success: false,
	// 		Message: "Internal server error",
	// 		Errors:  err.Error(),
	// 	})
	// 	return
	// }
	// c.JSON(http.StatusCreated, utils.Response{
	// 	Success: true,
	// 	Message: "Success to add new movie",
	// })
}

func UpdateMovie(c *gin.Context) {}

func DeleteMovie(c *gin.Context) {}

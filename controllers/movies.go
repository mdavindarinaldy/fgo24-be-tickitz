package controllers

import (
	"be-tickitz/dto"
	"be-tickitz/models"
	"be-tickitz/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	search := c.DefaultQuery("search", "")
	genre := c.DefaultQuery("genre", "")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	rdClient := utils.RedisConnect()
	endpoint := fmt.Sprintf("/movies?search=%s&genre=%s&page=%d", search, genre, page)
	pagepoint := fmt.Sprintf("/movies?search=%s&genre=%s&pagedata=%d", search, genre, page)
	result := rdClient.Exists(context.Background(), endpoint)
	if result.Val() == 0 {
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
		encodedMovies, _ := json.Marshal(movies)
		encodedPageData, _ := json.Marshal(pageData)
		rdClient.Set(context.Background(), endpoint, encodedMovies, 0)
		rdClient.Set(context.Background(), pagepoint, encodedPageData, 0)
	} else {
		data := rdClient.Get(context.Background(), endpoint)
		page := rdClient.Get(context.Background(), pagepoint)

		str := data.Val()
		movies := []dto.Movie{}
		pageData := page.Val()
		pageResult := []utils.PageData{}

		json.Unmarshal([]byte(str), &movies)
		json.Unmarshal([]byte(pageData), &pageResult)

		c.JSON(http.StatusOK, utils.Response{
			Success:  true,
			Message:  "Success to get movies",
			PageInfo: pageResult,
			Result:   movies,
		})
	}
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

	rdClient := utils.RedisConnect()
	endpoint := fmt.Sprintf("/movies/:%d", id)

	result := rdClient.Exists(context.Background(), endpoint)
	if result.Val() == 0 {
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
		encodedMovie, _ := json.Marshal(movie)
		rdClient.Set(context.Background(), endpoint, encodedMovie, 0)
	} else {
		data := rdClient.Get(context.Background(), endpoint)
		str := data.Val()
		movie := dto.Movie{}
		json.Unmarshal([]byte(str), &movie)
		c.JSON(http.StatusOK, utils.Response{
			Success: true,
			Message: "Success to get data",
			Result:  movie,
		})
	}
}

// GetUpcomingMovies retrieves a list of movies that has not been released yet
// @Summary Get upcoming movies
// @Description Retrieve a list of upcoming movies
// @Tags Movies
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response "Successful response with upcoming movies"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /movies/upcoming [get]
func GetUpcomingMovies(c *gin.Context) {
	rdClient := utils.RedisConnect()
	endpoint := "/movies/upcoming"
	result := rdClient.Exists(context.Background(), endpoint)

	if result.Val() == 0 {
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
		encodedMovies, _ := json.Marshal(movies)
		rdClient.Set(context.Background(), endpoint, encodedMovies, 0)
	} else {
		data := rdClient.Get(context.Background(), endpoint)
		str := data.Val()
		movies := []dto.Movie{}
		json.Unmarshal([]byte(str), &movies)
		c.JSON(http.StatusOK, utils.Response{
			Success: true,
			Message: "Success to get data",
			Result:  movies,
		})
	}
}

// GetGenres retrieves the list of all genres
// @Summary Get all genres
// @Description Retrieve a list of all genres (admin only)
// @Tags Admin: Movies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{result=[]dto.SubData} "Successful response with genres list"
// @Failure 401 {object} utils.Response "Unauthorized access"
// @Failure 403 {object} utils.Response "Forbidden access (requires admin role)"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /admin/genres [get]
func GetGenres(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, utils.Response{
			Success: false,
			Message: "Forbidden",
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
// @Tags Admin: Movies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param search query string false "Search by director name (case-insensitive)" example:"nolan"
// @Success 200 {object} utils.Response{result=[]dto.SubData} "Successful response with directors list"
// @Failure 401 {object} utils.Response "Unauthorized access"
// @Failure 403 {object} utils.Response "Forbidden access (requires admin role)"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /admin/directors [get]
func GetDirectors(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, utils.Response{
			Success: false,
			Message: "Forbidden",
		})
		return
	}
	search := c.DefaultQuery("search", "")
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
// @Tags Admin: Movies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param search query string false "Search by cast name (case-insensitive)" example:"dicaprio"
// @Success 200 {object} utils.Response{result=[]dto.SubData} "Successful response with casts list"
// @Failure 401 {object} utils.Response "Unauthorized access"
// @Failure 403 {object} utils.Response "Forbidden access (requires admin role)"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /admin/casts [get]
func GetCasts(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, utils.Response{
			Success: false,
			Message: "Forbidden",
		})
		return
	}
	search := c.DefaultQuery("search", "")
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

// AddDirectorHandler adds a new director
// @Summary Add a new director
// @Description Create a new director with name (admin only)
// @Tags Admin: Movies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param director body dto.NewData true "Director data"
// @Success 201 {object} utils.Response{result=dto.SubData} "Director created successfully"
// @Failure 400 {object} utils.Response{errors=string} "Bad request (e.g., empty director name)"
// @Failure 401 {object} utils.Response "Unauthorized access"
// @Failure 403 {object} utils.Response "Forbidden access (requires admin role)"
// @Failure 500 {object} utils.Response{errors=string} "Internal server error"
// @Router /admin/directors [post]
func AddDirector(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, utils.Response{
			Success: false,
			Message: "Forbidden",
		})
		return
	}

	newDirector := dto.NewData{}
	c.ShouldBind(&newDirector)
	data, err := models.AddDirector(newDirector.Name)

	if err != nil {
		if err.Error() == "director name should not be empty" {
			c.JSON(http.StatusInternalServerError, utils.Response{
				Success: false,
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Success to add new director",
		Result:  data,
	})
}

// AddCastHandler adds a new cast
// @Summary Add a new cast
// @Description Create a new cast with name (admin only)
// @Tags Admin: Movies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param cast body dto.NewData true "Cast data"
// @Success 201 {object} utils.Response{result=dto.SubData} "Cast created successfully"
// @Failure 400 {object} utils.Response{errors=string} "Bad request (e.g., empty cast name)"
// @Failure 401 {object} utils.Response "Unauthorized access"
// @Failure 403 {object} utils.Response "Forbidden access (requires admin role)"
// @Failure 500 {object} utils.Response{errors=string} "Internal server error"
// @Router /admin/casts [post]
func AddCast(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, utils.Response{
			Success: false,
			Message: "Forbidden",
		})
		return
	}

	newCast := dto.NewData{}
	c.ShouldBind(&newCast)
	data, err := models.AddCast(newCast.Name)

	if err != nil {
		if err.Error() == "cast name should not be empty" {
			c.JSON(http.StatusInternalServerError, utils.Response{
				Success: false,
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Success to add new cast",
		Result:  data,
	})
}

// AddGenreHandler adds a new genre
// @Summary Add a new genre
// @Description Create a new genre with name (admin only)
// @Tags Admin: Movies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param genre body dto.NewData true "Genre data"
// @Success 201 {object} utils.Response{result=dto.SubData} "Genre created successfully"
// @Failure 400 {object} utils.Response{errors=string} "Bad request (e.g., empty genre name)"
// @Failure 401 {object} utils.Response "Unauthorized access"
// @Failure 403 {object} utils.Response "Forbidden access (requires admin role)"
// @Failure 500 {object} utils.Response{errors=string} "Internal server error"
// @Router /admin/genres [post]
func AddGenre(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, utils.Response{
			Success: false,
			Message: "Forbidden",
		})
		return
	}

	newGenre := dto.NewData{}
	c.ShouldBind(&newGenre)
	data, err := models.AddGenre(newGenre.Name)

	if err != nil {
		if err.Error() == "genre name should not be empty" {
			c.JSON(http.StatusInternalServerError, utils.Response{
				Success: false,
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Success to add new genre",
		Result:  data,
	})
}

// AddMovie adds a new movie
// @Summary Add a new movie
// @Description Create a new movie with associated genres, directors, and casts (admin only)
// @Tags Admin: Movies
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param title formData string true "Movie title"
// @Param synopsis formData string true "Movie synopsis"
// @Param releaseDate formData string true "Release date (YYYY-MM-DD)"
// @Param price formData number true "Ticket price"
// @Param runtime formData int true "Duration in minutes"
// @Param genres formData string true "Comma-separated genre IDs"
// @Param directors formData string true "Comma-separated director IDs"
// @Param casts formData string true "Comma-separated cast IDs"
// @Param poster formData file true "Poster image"
// @Param backdrop formData file true "Backdrop image"
// @Success 201 {object} utils.Response "Movie created successfully"
// @Failure 400 {object} utils.Response "Bad request (e.g., empty movie data)"
// @Failure 401 {object} utils.Response "Unauthorized access"
// @Failure 403 {object} utils.Response "Forbidden access (requires admin role)"
// @Failure 500 {object} utils.Response{errors=string} "Internal server error"
// @Router /admin/movie [post]
func AddMovie(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, utils.Response{
			Success: false,
			Message: "Forbidden",
		})
		return
	}
	userId, _ := c.Get("userId")
	newMovie := dto.NewMovie{}
	err := c.ShouldBind(&newMovie)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Invalid input",
		})
		return
	}

	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}

	posterFile, _ := c.FormFile("poster")
	posterFileName := ""
	if posterFile != nil {
		if posterFile.Size > 5*1024*1024 {
			c.JSON(http.StatusBadRequest, utils.Response{
				Success: false,
				Message: "File is too large",
			})
			return
		}
		ext := strings.ToLower(filepath.Ext(posterFile.Filename))

		if !allowedExts[ext] {
			c.JSON(http.StatusBadRequest, utils.Response{
				Success: false,
				Message: "Invalid file type. Only JPG, JPEG, PNG allowed",
			})
			return
		}
		fileExt := filepath.Ext(posterFile.Filename)
		posterFileName = uuid.New().String() + fileExt
		err := c.SaveUploadedFile(posterFile, "./uploads/poster/"+posterFileName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, utils.Response{
				Success: false,
				Message: "Failed to save uploaded file",
			})
			return
		}
		newMovie.Poster = &posterFileName
	}

	backdrop, _ := c.FormFile("backdrop")
	backdropName := ""
	if backdrop != nil {
		if backdrop.Size > 5*1024*1024 {
			c.JSON(http.StatusBadRequest, utils.Response{
				Success: false,
				Message: "File is too large",
			})
			return
		}
		ext := strings.ToLower(filepath.Ext(backdrop.Filename))

		if !allowedExts[ext] {
			c.JSON(http.StatusBadRequest, utils.Response{
				Success: false,
				Message: "Invalid file type. Only JPG, JPEG, PNG allowed",
			})
			return
		}
		fileExt := filepath.Ext(backdrop.Filename)
		backdropName = uuid.New().String() + fileExt
		err := c.SaveUploadedFile(backdrop, "./uploads/backdrop/"+backdropName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, utils.Response{
				Success: false,
				Message: "Failed to save uploaded file",
			})
			return
		}
		newMovie.Backdrop = &backdropName
	}

	err = models.AddMovie(newMovie, int(userId.(float64)))
	if err != nil {
		if err.Error() == "new movie data should not be empty" {
			c.JSON(http.StatusBadRequest, utils.Response{
				Success: false,
				Message: err.Error(),
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
		Message: "Success to add new movie",
	})
	rdClient := utils.RedisConnect()
	keys := rdClient.Keys(context.Background(), "/movies?*").Val()
	for _, key := range keys {
		rdClient.Del(context.Background(), key)
	}
	rdClient.Del(context.Background(), "/movies/upcoming")
}

// UpdateMovieHandler updates an existing movie
// @Summary Update a movie
// @Description Update a movie's details and associated genres, directors, and casts (admin only)
// @Tags Admin: Movies
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param title formData string true "Movie title"
// @Param synopsis formData string true "Movie synopsis"
// @Param releaseDate formData string true "Release date (YYYY-MM-DD)"
// @Param price formData number true "Ticket price"
// @Param runtime formData int true "Duration in minutes"
// @Param genres formData string true "Comma-separated genre IDs"
// @Param directors formData string true "Comma-separated director IDs"
// @Param casts formData string true "Comma-separated cast IDs"
// @Param poster formData file true "Poster image"
// @Param backdrop formData file true "Backdrop image"
// @Success 200 {object} utils.Response "Movie updated successfully"
// @Failure 400 {object} utils.Response "Bad request (e.g., invalid input)"
// @Failure 401 {object} utils.Response "Unauthorized access"
// @Failure 403 {object} utils.Response "Forbidden access (requires admin role)"
// @Failure 500 {object} utils.Response{errors=string} "Internal server error"
// @Router /admin/movie/{id} [put]
func UpdateMovie(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, utils.Response{
			Success: false,
			Message: "Forbidden",
		})
		return
	}
	movieId, _ := strconv.Atoi(c.Param("id"))
	updateMovie := dto.NewMovie{}
	err := c.ShouldBind(&updateMovie)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Invalid input",
		})
		return
	}

	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}

	posterFile, _ := c.FormFile("poster")
	posterFileName := ""
	if posterFile != nil {
		if posterFile.Size > 5*1024*1024 {
			c.JSON(http.StatusBadRequest, utils.Response{
				Success: false,
				Message: "File is too large",
			})
			return
		}
		ext := strings.ToLower(filepath.Ext(posterFile.Filename))

		if !allowedExts[ext] {
			c.JSON(http.StatusBadRequest, utils.Response{
				Success: false,
				Message: "Invalid file type. Only JPG, JPEG, PNG allowed",
			})
			return
		}
		fileExt := filepath.Ext(posterFile.Filename)
		posterFileName = uuid.New().String() + fileExt
		err := c.SaveUploadedFile(posterFile, "./uploads/poster/"+posterFileName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, utils.Response{
				Success: false,
				Message: "Failed to save uploaded file",
			})
			return
		}
		updateMovie.Poster = &posterFileName
	}

	backdrop, _ := c.FormFile("backdrop")
	backdropName := ""
	if backdrop != nil {
		if backdrop.Size > 5*1024*1024 {
			c.JSON(http.StatusBadRequest, utils.Response{
				Success: false,
				Message: "File is too large",
			})
			return
		}
		ext := strings.ToLower(filepath.Ext(backdrop.Filename))

		if !allowedExts[ext] {
			c.JSON(http.StatusBadRequest, utils.Response{
				Success: false,
				Message: "Invalid file type. Only JPG, JPEG, PNG allowed",
			})
			return
		}
		fileExt := filepath.Ext(backdrop.Filename)
		backdropName = uuid.New().String() + fileExt
		err := c.SaveUploadedFile(backdrop, "./uploads/backdrop/"+backdropName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, utils.Response{
				Success: false,
				Message: "Failed to save uploaded file",
			})
			return
		}
		updateMovie.Backdrop = &backdropName
	}

	err = models.UpdateMovie(updateMovie, movieId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Update movie data success",
	})
	rdClient := utils.RedisConnect()
	keys := rdClient.Keys(context.Background(), "/movies?*").Val()
	for _, key := range keys {
		rdClient.Del(context.Background(), key)
	}
	rdClient.Del(context.Background(), "/movies/upcoming")
}

// DeleteMovie deletes a movie
// @Summary Delete a movie
// @Description Delete a movie by ID, including its associated genres, directors, and casts (admin only)
// @Tags Admin: Movies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Movie ID"
// @Success 200 {object} utils.Response "Movie deleted successfully"
// @Failure 400 {object} utils.Response "Bad request (e.g., invalid movie ID)"
// @Failure 401 {object} utils.Response "Unauthorized access"
// @Failure 403 {object} utils.Response "Forbidden access (requires admin role)"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /admin/movie/{id} [delete]
func DeleteMovie(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, utils.Response{
			Success: false,
			Message: "Forbidden",
		})
		return
	}

	movieId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Invalid id",
		})
		return
	}

	err = models.DeleteMovie(movieId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to delete movie",
		})
		return
	}

	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Delete movie success",
	})
	rdClient := utils.RedisConnect()
	keys := rdClient.Keys(context.Background(), "/movies?*").Val()
	for _, key := range keys {
		rdClient.Del(context.Background(), key)
	}
	rdClient.Del(context.Background(), "/movies/upcoming")
	rdClient.Del(context.Background(), fmt.Sprintf("/movies/upcoming/%d", movieId))
}

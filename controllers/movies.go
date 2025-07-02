package controllers

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetMovies(c *gin.Context) {
	search := strings.ToLower(c.DefaultQuery("search", ""))
	genre := strings.ToLower(c.DefaultQuery("genre", ""))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
}

func GetUpcomingMovies(c *gin.Context) {

}

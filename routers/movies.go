package routers

import (
	"be-tickitz/controllers"

	"github.com/gin-gonic/gin"
)

func moviesRouter(r *gin.RouterGroup) {
	r.GET("", controllers.GetMovies)
	r.GET("/:id", controllers.GetDetailMovie)
	r.GET("/upcoming", controllers.GetUpcomingMovies)
	r.GET("/genres", controllers.GetGenres)
	r.GET("/directors", controllers.GetDirectors)
	r.GET("/casts", controllers.GetCasts)
}

package routers

import (
	"be-tickitz/controllers"

	"github.com/gin-gonic/gin"
)

func moviesRouter(r *gin.RouterGroup) {
	r.GET("", controllers.GetMovies)
	r.GET("/:id", controllers.GetDetailMovie)
	r.GET("/upcoming", controllers.GetUpcomingMovies)
}

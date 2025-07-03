package routers

import (
	"be-tickitz/controllers"
	"be-tickitz/middlewares"

	"github.com/gin-gonic/gin"
)

func moviesRouter(r *gin.RouterGroup) {
	r.GET("", controllers.GetMovies)
	r.GET("/:id", controllers.GetDetailMovie)
	r.GET("/upcoming", controllers.GetUpcomingMovies)
	r.GET("/genres", middlewares.VerifyToken(), controllers.GetGenres)
	r.GET("/directors", middlewares.VerifyToken(), controllers.GetDirectors)
	r.GET("/casts", middlewares.VerifyToken(), controllers.GetCasts)
	r.POST("/genres", middlewares.VerifyToken(), controllers.AddGenre)
	r.POST("/directors", middlewares.VerifyToken(), controllers.AddDirector)
	r.POST("/casts", middlewares.VerifyToken(), controllers.AddCast)
	r.POST("", middlewares.VerifyToken(), controllers.AddMovie)
	r.PUT("/:id", middlewares.VerifyToken(), controllers.UpdateMovie)
	r.DELETE("/:id", middlewares.VerifyToken(), controllers.DeleteMovie)
}

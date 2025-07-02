package routers

import (
	"be-tickitz/controllers"

	"github.com/gin-gonic/gin"
)

func moviesRouter(r *gin.RouterGroup) {
	r.GET("", controllers.GetMovies)
}

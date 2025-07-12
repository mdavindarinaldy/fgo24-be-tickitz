package routers

import (
	"be-tickitz/controllers"
	"be-tickitz/middlewares"

	"github.com/gin-gonic/gin"
)

func adminRouter(r *gin.RouterGroup) {
	r.Use(middlewares.VerifyToken())
	r.GET("/genres", controllers.GetGenres)
	r.GET("/directors", controllers.GetDirectors)
	r.GET("/casts", controllers.GetCasts)
	r.POST("/genres", controllers.AddGenre)
	r.POST("/directors", controllers.AddDirector)
	r.POST("/casts", controllers.AddCast)
	r.POST("/movie", controllers.AddMovie)
	r.PUT("/movie/:id", controllers.UpdateMovie)
	r.DELETE("/movie/:id", controllers.DeleteMovie)
	r.POST("/payment-methods", controllers.AddPaymentMethod)
	r.GET("/sales", controllers.GetSalesPerMovie)
}

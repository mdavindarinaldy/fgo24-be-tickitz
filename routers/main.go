package routers

import (
	"be-tickitz/docs"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func CombineRouter(r *gin.Engine) {
	authRouter(r.Group("/auth"))
	moviesRouter(r.Group("/movies"))
	adminRouter(r.Group("/admin"))
	transactionsRouter(r.Group("/transactions"))
	userRouter(r.Group("/profile"))
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/docs", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusPermanentRedirect, "/docs/index.html")
	})
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

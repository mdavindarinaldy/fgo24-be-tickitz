package routers

import (
	"be-tickitz/controllers"

	"github.com/gin-gonic/gin"
)

func authRouter(r *gin.RouterGroup) {
	r.POST("/register", controllers.AuthRegister)
}

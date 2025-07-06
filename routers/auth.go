package routers

import (
	"be-tickitz/controllers"

	"github.com/gin-gonic/gin"
)

func authRouter(r *gin.RouterGroup) {
	r.POST("/register", controllers.AuthRegister)
	r.POST("/login", controllers.AuthLogin)
	r.POST("/pass", controllers.AuthForgotPass)
	r.PATCH("/pass", controllers.AuthResetPass)
	r.POST("/logout", controllers.AuthLogout)
}

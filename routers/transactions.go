package routers

import (
	"be-tickitz/controllers"
	"be-tickitz/middlewares"

	"github.com/gin-gonic/gin"
)

func transactionsRouter(r *gin.RouterGroup) {
	r.POST("/payment-method", middlewares.VerifyToken(), controllers.AddPaymentMethod)
}

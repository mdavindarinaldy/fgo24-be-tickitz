package routers

import (
	"be-tickitz/controllers"
	"be-tickitz/middlewares"

	"github.com/gin-gonic/gin"
)

func transactionsRouter(r *gin.RouterGroup) {
	r.Use(middlewares.VerifyToken())
	r.POST("/payment-methods", controllers.AddPaymentMethod)
	r.GET("/payment-methods", controllers.GetPaymentMethod)
	r.POST("", controllers.AddTransactions)
	r.GET("", controllers.GetTransactionsHistory)
	r.GET("/seats", controllers.GetReservedSeat)
}

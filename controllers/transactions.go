package controllers

import (
	"be-tickitz/dto"
	"be-tickitz/models"
	"be-tickitz/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func BookTicket() {

}

// AddPaymentMethod adds a new payment method
// @Summary Add a new payment method
// @Description Create a new payment method (admin only)
// @Tags Transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param movie body dto.NewData true "New payment method data"
// @Success 201 {object} utils.Response "Payment method created successfully"
// @Failure 400 {object} utils.Response "Bad request (e.g., empty payment method data)"
// @Failure 401 {object} utils.Response "Unauthorized access (requires admin role)"
// @Failure 500 {object} utils.Response{errors=string} "Internal server error"
// @Router /transactions/payment-method [post]
func AddPaymentMethod(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusUnauthorized, utils.Response{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}
	newMethod := dto.NewData{}
	c.ShouldBind(&newMethod)
	err := models.AddPaymentMethod(newMethod)
	if err != nil {
		if err.Error() == "payment method name should not be empty" {
			c.JSON(http.StatusBadRequest, utils.Response{
				Success: false,
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Status internal server error",
		})
		return
	}
	c.JSON(http.StatusCreated, utils.Response{
		Success: true,
		Message: "Payment method successfully been added",
	})
}

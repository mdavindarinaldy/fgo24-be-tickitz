package controllers

import (
	"be-tickitz/dto"
	"be-tickitz/models"
	"be-tickitz/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AddPaymentMethod adds a new payment method
// @Summary Add a new payment method
// @Description Create a new payment method (admin only)
// @Tags Admin: Transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param movie body dto.NewData true "New payment method data"
// @Success 201 {object} utils.Response "Payment method created successfully"
// @Failure 400 {object} utils.Response "Bad request (e.g., empty payment method data)"
// @Failure 401 {object} utils.Response "Unauthorized access"
// @Failure 403 {object} utils.Response "Forbidden access (requires admin role)"
// @Failure 500 {object} utils.Response{errors=string} "Internal server error"
// @Router /admin/payment-methods [post]
func AddPaymentMethod(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, utils.Response{
			Success: false,
			Message: "Forbidden",
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

// GetPaymentMethod retrieves payment method
// @Summary Get payment method
// @Description Retrieve a list of payment method
// @Tags Transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{result=[]dto.SubData} "Successful response with payment methods list"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /transactions/payment-methods [get]
func GetPaymentMethod(c *gin.Context) {
	data, err := models.GetPaymentMethod()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to get data",
		})
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Success to get data",
		Result:  data,
	})
}

// AddTransactions creates a new transaction for booking movie tickets
// @Summary Create a new transaction
// @Description Creates a transaction for booking movie tickets, including checking or creating a showtime and reserving seats (user only)
// @Tags Transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param newTrx body dto.NewTrx true "Transaction data"
// @Success 201 {object} utils.Response{result=dto.TrxSuccess} "Transaction created successfully"
// @Failure 400 {object} utils.Response{errors=string} "Bad request due to invalid input"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /transactions [post]
func AddTransactions(c *gin.Context) {
	userId, _ := c.Get("userId")
	newTrx := dto.NewTrx{}
	c.ShouldBind(&newTrx)
	showtimeId, transactionId, err := models.AddTransactions(newTrx, int(userId.(float64)))
	if err != nil {
		if err.Error() == "transactions data should not be empty" || strings.Contains(err.Error(), "unique_seat_per_showtime") || err.Error() == "invalid showtime format, use HH:MM:SS" || err.Error() == "invalid date format, use YYYY-MM-DD" {
			c.JSON(http.StatusBadRequest, utils.Response{
				Success: false,
				Message: "Failed to order ticket",
				Errors:  err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Internal server error",
			Errors:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, utils.Response{
		Success: true,
		Message: "Success to order ticket",
		Result: dto.TrxSuccess{
			ShowtimeId:    showtimeId,
			TransactionId: transactionId,
		},
	})
	rdClient := utils.RedisConnect()
	rdClient.Del(context.Background(), fmt.Sprintf("/transactions/:%d", int(userId.(float64))))
}

// GetTransactionsHistory retrieves list of user's transactions
// @Summary Get transactions history
// @Description Retrieve a list of user's transactions
// @Tags Transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{result=[]dto.TransactionHistory} "Successful response with transactions history list"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /transactions [get]
func GetTransactionsHistory(c *gin.Context) {
	userId, _ := c.Get("userId")
	rdClient := utils.RedisConnect()
	endpoint := fmt.Sprintf("/transactions/:%d", int(userId.(float64)))
	result := rdClient.Exists(context.Background(), endpoint)
	if result.Val() == 0 {
		trxHistory, err := models.GetTransactionsHistory(int(userId.(float64)))
		if err != nil {
			c.JSON(http.StatusInternalServerError, utils.Response{
				Success: false,
				Message: "Failed to get data",
				Errors:  err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, utils.Response{
			Success: true,
			Message: "Success to get data",
			Result:  trxHistory,
		})
		endcodedTrx, _ := json.Marshal(trxHistory)
		rdClient.Set(context.Background(), endpoint, endcodedTrx, 0)
	} else {
		data := rdClient.Get(context.Background(), endpoint)
		str := data.Val()
		trxHistory := []dto.TransactionHistory{}
		json.Unmarshal([]byte(str), &trxHistory)
		c.JSON(http.StatusOK, utils.Response{
			Success: true,
			Message: "Success to get movies",
			Result:  trxHistory,
		})
	}
}

// GetReservedSeat retrieves reserved seats for a showtime
// @Summary Get reserved seats
// @Description Retrieves a list of reserved seats for a specific showtime identified by movie ID, cinema, location, date, and showtime
// @Tags Transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id_movie query int true "Movie ID"
// @Param cinema query string true "Cinema name"
// @Param location query string true "Cinema location"
// @Param date query string true "Showtime date (YYYY-MM-DD)"
// @Param showtime query string true "Showtime (HH:MM:SS)"
// @Success 200 {object} utils.Response{result=dto.ReservedSeatsResponse} "Successful response with reserved seats"
// @Failure 400 {object} utils.Response{errors=string} "Bad request due to invalid input"
// @Failure 500 {object} utils.Response{errors=string} "Internal server error"
// @Router /transactions/seats [get]
func GetReservedSeat(c *gin.Context) {
	var req dto.ReservedSeatsRequest
	c.ShouldBind(&req)
	res, err := models.GetReservedSeat(req)
	if err != nil {
		if err.Error() == "all fields must be provided" {
			c.JSON(http.StatusBadRequest, utils.Response{
				Success: false,
				Message: "Failed to get data reserved seats",
				Errors:  err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to get data reserved seats",
		})
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Success to get data reserved seats",
		Result:  res,
	})
}

// GetSalesPerMovie retrieves total sales data per movie
// @Summary Get sales data per movie
// @Description Retrieves aggregated sales data for each movie, accessible only to admin users
// @Tags Admin: Transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 201 {object} utils.Response{result=[]dto.SalesPerMovie} "Successful response with sales data per movie"
// @Failure 401 {object} utils.Response "Unauthorized access"
// @Failure 403 {object} utils.Response "Forbidden access (requires admin role)"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /admin/sales [get]
func GetSalesPerMovie(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, utils.Response{
			Success: false,
			Message: "Forbidden",
		})
		return
	}
	result, err := models.GetSalesPerMovie()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Status internal server error",
			Errors:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, utils.Response{
		Success: true,
		Message: "Success to get sales data",
		Result:  result,
	})
}

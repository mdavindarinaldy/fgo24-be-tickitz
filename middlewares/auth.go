package middlewares

import (
	"be-tickitz/utils"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		secretKey := os.Getenv("APP_SECRET")
		token := strings.Split(c.GetHeader("Authorization"), "Bearer ")

		if len(token) < 2 {
			c.JSON(http.StatusUnauthorized, utils.Response{
				Success: false,
				Message: "Unauthorized",
			})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		rawToken, err := jwt.Parse(token[1], func(t *jwt.Token) (any, error) {
			return []byte(secretKey), nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.Response{
				Success: false,
				Message: "Invalid token",
			})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		userId := rawToken.Claims.(jwt.MapClaims)["userId"]
		role := rawToken.Claims.(jwt.MapClaims)["role"]
		c.Set("userId", userId)
		c.Set("role", role)
		c.Next()
	}
}
